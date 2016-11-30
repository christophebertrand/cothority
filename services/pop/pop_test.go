package pop

import (
	"testing"
	"time"
	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/sda"
	"github.com/dedis/cothority/network"
	"github.com/dedis/crypto/abstract"

)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func NewTestClient(lt *sda.LocalTest) *Client {
	return &Client{Client: lt.NewClient(ServiceName)}
}

//Sets up an example configuration file
func setupConfigFile() *ConfigurationFile{
	rand := network.Suite.Cipher([]byte("example"))

	X := make([]abstract.Point, 3)
	for i := range X { // pick random points
		x := network.Suite.Scalar().Pick(rand) // create a private key x
    	X[i] = network.Suite.Point().Mul(nil, x)
	}
	return &ConfigurationFile{
		OrganizersPublic : X,
		StartingTime: 5.5,
		EndingTime: 5.8,
		Duration: 66.6,
		Context: []byte("IFF Forum"),
		Date: time.Date(2016, 5, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestSendConfigFileHash(t *testing.T) {
	local := sda.NewLocalTest()
	// generate 5 hosts, they don't connect, they process messages, and they
	// don't register the tree or entitylist
	_, el, _ := local.GenTree(5, true)
	defer local.CloseAll()
	dst := el.RandomServerIdentity() //For now a random server
	client := NewTestClient(local)
	config := setupConfigFile()
	hash_value, err := client.SendConfigFileHash(dst,config)
	log.ErrFatal(err, "Problem inside SendConfigFileHash")
	log.Lvl1("Config File was hashed with ",hash_value)
}
