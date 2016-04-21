package manage

import (
	"errors"

	"github.com/dedis/cothority/lib/dbg"
	"github.com/dedis/cothority/lib/sda"
)

func init() {
	sda.ProtocolRegisterName("Broadcast", NewBroadcastProtocol)
}

// Broadcast ensures that all nodes are connected to each other. If you need
// a confirmation once everything is set up, you can register a callback-function
// using RegisterOnDone()
type Broadcast struct {
	*sda.TreeNodeInstance
	onDoneCb    func()
	repliesLeft int
	tnIndex     int
}

// NewBroadcastProtocol returns a new Broadcast protocol
func NewBroadcastProtocol(n *sda.TreeNodeInstance) (sda.ProtocolInstance, error) {
	b := &Broadcast{
		TreeNodeInstance: n,
		tnIndex:          -1,
	}
	for i, tn := range n.Tree().List() {
		if tn.Id == n.TreeNode().Id {
			b.tnIndex = i
		}
	}
	if b.tnIndex == -1 {
		return nil, errors.New("Didn't find my TreeNode in the Tree")
	}
	err := n.RegisterHandler(b.handleContactNodes)
	if err != nil {
		return nil, err
	}
	err = n.RegisterHandler(b.handleDone)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Start will contact everyone and makes the connections
func (b *Broadcast) Start() error {
	n := len(b.Tree().List())
	b.repliesLeft = n * (n - 1) / 2
	b.SendTo(b.Root(), &ContactNodes{})
	dbg.Lvl3(b.Name(), "Sent Announce to everyone")
	return nil
}

// handleAnnounce receive the announcement from another node
// it reply with an ACK.
func (b *Broadcast) handleContactNodes(msg struct {
	*sda.TreeNode
	ContactNodes
}) {
	dbg.Lvl3(b.Info(), "Received message from", msg.TreeNode.String())
	if msg.TreeNode.Id == b.Root().Id {
		dbg.Lvl3(b.Info(), "Contacting everybody")
		// Connect to all nodes that are later in the TreeNodeList, but only if
		// the message comes from root
		for _, tn := range b.Tree().List()[b.tnIndex+1:] {
			dbg.Lvl3("Connecting to", tn.String())
			err := b.SendTo(tn, &ContactNodes{})
			if err != nil {
				return
			}
		}
	}
	// Tell Root we're done
	b.SendTo(b.Root(), &Done{})
}

// Every node being contacted sends back a Done to the root which has
// to count to decide if all is done
func (b *Broadcast) handleDone(struct {
	*sda.TreeNode
	Done
}) {
	b.repliesLeft--
	dbg.Lvl3("Got reply and waiting for more:", b.repliesLeft)
	if b.repliesLeft == 0 {
		if b.onDoneCb != nil {
			dbg.Lvl2("Done with broadcasting to everybody")
			b.onDoneCb()
		}
	}
}

// RegisterOnDone takes a function that will be called once all connections
// are set up.
func (b *Broadcast) RegisterOnDone(fn func()) {
	b.onDoneCb = fn
}

// ContactNodes is sent from the root to ALL other nodes
type ContactNodes struct{}

// Done is sent back to root once everybody has been contacted
type Done struct{}
