package gwent

// IaMove returns the card played by the IA, or nil if it passes
func IaMove(game *Game, ia *GameSide, enemy *GameSide) (*Card, Row) {
	all_board := game.Merge()
	if enemy.Passed && ia.CachedScore > enemy.CachedScore {
		return nil, 0 // pass
	}
	if all_board.Len() == 0 {
		//IA is starting the round
	}
	if game.Round() == 0 { //at first round, play all the spies
		if spy := ia.Hand.GetByEffect(Spy); spy != nil {
			return spy, spy.Row
		}
	}
	card := ia.Hand.GetByType(TypeNormal)
	if card != nil {
		return card, IaChooseRow(card, ia)

	}
	//no unit card found, IA will play weather & scorch
	if card := BestMove(game); card != nil {
		return card, card.Row
	}
	return card, IaChooseRow(card, ia)
}

func IaChooseRow(card *Card, ia *GameSide) Row {
	if card == nil {
		return 0
	}
	if card.Effects.Has(Agile) {
		if ia.ScoreRangedCombat > ia.ScoreCloseCombat {
			return RangedCombat
		} else {
			return CloseCombat
		}
	}
	return card.Row
}

// TestSim returns the score after the simulated move
func TestSim(game *Game, card *Card, row Row) (player int, enemy int) {
	sim := game.DeepCopy()
	if sim == nil {
		println(game.Turn.String())
		panic("deepcopy error")
	}

	sim_card := sim.Player().Hand.GetById(card.Id) //sim_hand.GetById(card.Id)
	if sim_card == nil {
		panic("error in sim")
	}

	sim.PlayMove(sim_card, Special, sim.Player(), sim.Enemy())
	return sim.Player().CachedScore, sim.Enemy().CachedScore
}

//BestMove finds the card in hand that maximizes the difference of score between foes (positive = IA wins)
// if the best move ends up worse than doing nothing, returns nil
func BestMove(game *Game) *Card {
	differences := map[int]*Card{} //score->card
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

type CardType enum

const (
	TypeNormal  CardType = iota // basic, heroes, spies, healers
	TypeWeather                 // rain/frost/fog/sun
	TypeSpecial                 // scorch, morale boost
)

func (card *Card) Type() CardType {
	if card.Row == Weather {
		return TypeWeather
	}
	if card.Effects.Has(Scorch) && card.Strength == 0 {
		return TypeSpecial
	}
	return TypeNormal
}

func (hand *CardList) GetByEffect(effect Effect) *Card {
	for _, card := range *hand {
		if card.Effects.Has(effect) {
			return card
		}
	}
	return nil
}

func (hand *CardList) GetByType(cardtype CardType) *Card {
	for _, card := range *hand {
		if card.Type() == cardtype {
			return card
		}
	}
	return nil
}
