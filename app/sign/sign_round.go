package main

import (
	"github.com/dedis/cothority/lib/dbg"
	"github.com/dedis/cothority/lib/monitor"
	"github.com/dedis/cothority/lib/sign"
)

const RoundMeasureType = "measure"

type RoundMeasure struct {
	measure    *monitor.Measure
	firstRound int
	*sign.RoundCosi
}

// Pass firstround, as we will have some previous rounds to wait
// for everyone to be setup
func RegisterRoundMeasure(firstRound int) {
	sign.RegisterRoundFactory(RoundMeasureType,
		func(s *sign.Node) sign.Round {
			return NewRoundMeasure(s, firstRound)
		})
}

func NewRoundMeasure(node *sign.Node, firstRound int) *RoundMeasure {
	dbg.Lvlf3("Making new roundmeasure %+v", node)
	round := &RoundMeasure{}
	round.RoundCosi = sign.NewRoundCosi(node)
	round.Type = RoundMeasureType
	round.firstRound = firstRound
	return round
}

func (round *RoundMeasure) Announcement(viewNbr, roundNbr int, in *sign.AnnouncementMessage, out []*sign.AnnouncementMessage) error {
	if round.IsRoot {
		round.measure = monitor.NewMeasure("round")
	}
	return round.RoundCosi.Announcement(viewNbr, roundNbr, in, out)
}

func (round *RoundMeasure) Response(in []*sign.ResponseMessage, out *sign.ResponseMessage) error {
	err := round.RoundCosi.Response(in, out)
	if round.IsRoot {
		round.measure.Measure()
		dbg.Lvl1("Round", round.RoundNbr-round.firstRound+1,
			"finished - took", round.measure.WallTime)
	}
	return err
}
