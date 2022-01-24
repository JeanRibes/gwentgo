package ia

import (
	"gwentgo/gwent"
)

type GwyntAI interface {
	ChooseCard(game *gwent.Game, ia *gwent.GameSide, enemy *gwent.GameSide) (*gwent.Card, gwent.Row)
}

type BaseIA struct{}

func Init() *BaseIA {
	return &BaseIA{}
}

// IaMove returns the card played by the IA, or nil if it passes
func (this *BaseIA) ChooseCard(game *gwent.Game, ia *gwent.GameSide, enemy *gwent.GameSide) (*gwent.Card, gwent.Row) {
	if enemy.Passed && ia.CachedScore > enemy.CachedScore {
		return nil, 0 // pass
	}
	if game.Merge().Len() == 0 {
		//IA is starting the round
	}
	if game.Round() == 0 { //at first round, play all the spies
		if spy := ia.Hand.GetByEffect(gwent.Spy); spy != nil {
			return spy, spy.Row
		}
	}
	card := ia.Hand.GetByType(gwent.TypeNormal)
	if card != nil {
		return card, IaChooseRow(card, ia)

	}
	//no unit card found, IA will play weather & scorch
	if card := BestMove(game); card != nil {
		return card, card.Row
	}
	return card, IaChooseRow(card, ia)
}

func IaMove(game *gwent.Game, ia *gwent.GameSide, enemy *gwent.GameSide) (*gwent.Card, gwent.Row) {
	return Init().ChooseCard(game, ia, enemy)
}

func IaChooseRow(card *gwent.Card, ia *gwent.GameSide) gwent.Row {
	if card == nil {
		return 0
	}
	if card.Effects.Has(gwent.Agile) {
		if ia.ScoreRangedCombat > ia.ScoreCloseCombat {
			return gwent.RangedCombat
		} else {
			return gwent.CloseCombat
		}
	}
	return card.Row
}

// TestSim returns the score after the simulated move
func TestSim(game *gwent.Game, card *gwent.Card, row gwent.Row) (player int, enemy int) {
	sim := game.DeepCopy()
	if sim == nil {
		println(game.Turn.String())
		panic("deepcopy error")
	}

	sim_card := sim.Player().Hand.GetById(card.Id) //sim_hand.GetById(card.Id)
	if sim_card == nil {
		panic("error in sim")
	}

	sim.PlayMove(sim_card, gwent.Special, sim.Player(), sim.Enemy())
	return sim.Player().CachedScore, sim.Enemy().CachedScore
}

//BestMove finds the card in hand that maximizes the difference of score between foes (positive = IA wins)
// if the best move ends up worse than doing nothing, returns nil
func BestMove(game *gwent.Game) *gwent.Card {
	differences := map[int]*gwent.Card{} //score->card
	for _, card := range *game.Player().Hand {
		player, enemy := TestSim(game, card, card.Row)
		differences[player-enemy] = card
	}
	max := 0
	for diff, _ := range differences {
		if diff > max {
			max = diff
		}
	}
	for diff, card := range differences {
		if diff == max {
			if diff > (game.Player().CachedScore - game.Enemy().CachedScore) {
				// we dont want to make things worse by playing the least terribe card
				return card
			}
		}
	}
	return nil
}
