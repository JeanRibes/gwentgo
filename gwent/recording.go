package gwent

import (
	"bytes"
	"fmt"
)

func (g *Game) record(recordType RecordType, side Turn, card *Card, row Row) {
	rec := Record{
		Type: recordType,
		Side: side,
		Card: Card{},
		Row:  Special,
	}
	if side != Tie {
		p := g.Side(side)
		e := g.Side(g.EnemySide(side))
		rec.PlayerScore = RecordScore{
			Total:        p.CachedScore,
			CloseCombat:  p.ScoreCloseCombat,
			RangedCombat: p.ScoreRangedCombat,
			Siege:        p.ScoreSiege,
		}
		rec.EnemyScore = RecordScore{
			Total:        e.CachedScore,
			CloseCombat:  e.ScoreCloseCombat,
			RangedCombat: e.ScoreRangedCombat,
			Siege:        e.ScoreSiege,
		}
	}
	if card != nil {
		rec.Card = *card
		rec.Row = row
	}
	g.Recording = append(g.Recording, rec)
}

func (g *Game) recordMove(side Turn, card *Card, row Row) {
	g.record(CardMove, side, card, row)
}
func (g *Game) recordPass(side Turn) {
	g.record(Pass, side, nil, 0)
}
func (g *Game) recordRound(side Turn) {
	g.record(RoundEnd, side, nil, 0)
}

type RecordType enum

const (
	CardMove RecordType = iota
	RoundEnd            // encodes victory
	Pass                //encode the passing player
)

func (rt RecordType) String() string {
	if rt == CardMove {
		return "CardMove"
	}
	if rt == RoundEnd {
		return "RoundEnd"
	}
	if rt == Pass {
		return "Pass"
	}
	return ""
}

func (rt *RecordType) UnString(str string) {
	if str == "CardMove" {
		*rt = CardMove
	}
	if str == "RoundEnd" {
		*rt = RoundEnd
	}
	if str == "Pass" {
		*rt = Pass
	}
}

func (rt RecordType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("")
	buffer.WriteRune('"')
	buffer.WriteString(rt.String())
	buffer.WriteRune('"')
	return buffer.Bytes(), nil
}

type RecordScore struct {
	Total        int `json:"total"`
	CloseCombat  int `json:"close_combat"`
	RangedCombat int `json:"ranged_combat"`
	Siege        int `json:"siege"`
}

type Record struct {
	Type RecordType `json:"type"`
	Side Turn       `json:"side"`
	Card Card       `json:"card"`
	Row  Row        `json:"row"`

	PlayerScore RecordScore `json:"player_score"`
	EnemyScore  RecordScore `json:"enemy_score"`
}

func (rec Record) String() string {
	if rec.Type == CardMove {
		return fmt.Sprintf("%s %s %s %s", rec.Type.String(), rec.Side.String(), rec.Row, rec.Card.String())
	} else {
		return fmt.Sprintf("%s %s", rec.Type.String(), rec.Side.String())
	}
}
