package gwent

import (
	"math/rand"
)

func (game *Game) PlayMove(card *Card, slotRow Row, player *GameSide, enemy *GameSide) {
	if CheckMove(card, slotRow, player, enemy, game.WeatherCards) {
		game.WeatherClean() //needed at the beginning, otherwise debuff may apply in presence of the Clear weather

		if card.Effects.Has(Scorch) {
			game.Scorch()
			return
		}

		var row *Cards

		if slotRow == Weather {
			row = game.WeatherCards
		} else {
			if card.Effects.Has(Spy) {
				row = enemy.GetRow(slotRow)
			} else {
				row = player.GetRow(slotRow)
			}
		}
		MoveCard(player.Hand, row, card)

		/*post-move*/
		if card.Effects.Has(Spy) {
			player.Hand.Draw(player.Deck)
			player.Hand.Draw(player.Deck)
			//draw 2 from deck
		}
		if card.Effects.Has(Medic) {
			//draw from heap
			if player.Heap.Len() > 0 {
				player.Hand.Draw(player.Heap)
			}
		}
		if card.Effects.Has(Muster) {
			//find others in deck
			buddies := append(player.Hand.FindByName(card.Name), player.Deck.FindByName(card.Name)...)
			MoveCards(player.Deck, row, buddies)
		}
	}
}

func CheckMove(
	card *Card,
	slotRow Row,
	player *GameSide,
	enemy *GameSide,
	weathers *Cards,
) (valid bool) {
	if slotRow != card.Row {
		if card.Effects.Has(Agile) {
			card.Row = slotRow /* MUTATION */
		} else {
			return false
		}
	}
	return player.Hand.Has(card)
}

func InitHand(deck *Cards) (hand *Cards) {
	hand = &Cards{}
	for hand.Len() < 10 {
		hand.Draw(deck)
	}
	return hand
}

func (hand *Cards) Draw(deck *Cards) *Card {
	i := rand.Intn(len(*deck))
	for _, card := range *deck {
		i -= 1
		if i == 0 {
			MoveCard(deck, hand, card)
			return card
		}
	}
	panic(nil)
}

func InitDeck(all *[]Card, faction Faction) *Cards {
	deck := &Cards{}
	_all := *all
	for id, _ := range _all {
		card := _all[id]
		if card.Faction == faction || card.Faction == Neutral {
			card.Id = id
			card.Faction = faction
			deck.Add(&card)
		}
	}
	return deck
}

func (player *GameSide) Score(weather Effects) (sum int) {
	sum = 0
	sum += player.CloseCombat.Score(weather.Has(BitingFrost))
	sum += player.RangedCombat.Score(weather.Has(ImpenetrableFog))
	sum += player.Siege.Score(weather.Has(TorrentialRain))
	return sum
}

func (row Cards) Score(weatherDebuff bool) (sum int) {
	effects := row.Effects()
	moraleBoost := effects.Has(MoraleBoost)
	hornBoost := effects.Has(CommanderHorn)
	tightBonds := BondMap{}

	tightBonded := []*Card{}

	for _, card := range row {
		card.score = card.Score(weatherDebuff, moraleBoost, hornBoost)
		if card.Effects.Has(TightBond) {
			tightBonds.Add(card)
			tightBonded = append(tightBonded, card)
		}
	}

	for _, card := range tightBonded {
		/*card.score *= tightBonds.Get(card)*/
		if tightBonds.Get(card) >= 2 {
			card.score *= 2
		}
	}
	sum = 0
	for _, card := range row {
		sum += card.score
	}
	return sum
}

func (card *Card) Score(weatherDebuff bool, moraleBoost bool, hornBoost bool) (score int) {
	score = card.Strength
	if card.Effects.Has(Hero) {
		return
	}
	if weatherDebuff {
		score = 1
	}
	if hornBoost {
		score *= 2
	}
	if moraleBoost {
		score += 1
	}
	//log.Printf("score:%d for card %s", score, card)
	return
}

/*
 */
func (game *Game) Scorch() {
	weather := game.WeatherCards.Effects()
	//calculate score with side effects
	game.SideA.Score(weather)
	game.SideB.Score(weather)

	cards := game.Merge()
	maxScore := 0
	for _, card := range cards {
		if card.score > maxScore {
			maxScore = card.score
		}
	}
	for _, card := range cards {
		if card.score == maxScore {
			card.Strength = 0
			card.score = -1
		}
	}
	game.Clean() // deletes the scorched cards
}
