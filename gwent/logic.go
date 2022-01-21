package gwent

import (
	"log"
	"math/rand"
)

const MAX_ROUND = 3

func (game *Game) PlayMove(
	card *Card,
	slotRow Row,
	player *GameSide,
	enemy *GameSide) bool { // boolean: indicates if there is a "medic" choice
	if player.Hand.Has(card) && !player.Passed {

		game.WeatherClean() //needed at the beginning, otherwise debuff may apply in presence of the Clear weather

		if CheckMove(card, slotRow, player, enemy, game.WeatherCards) {
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
					//player.Hand.Draw(player.Heap)
					return true
				}
			}
			if card.Effects.Has(Muster) {
				//find others in deck and hand, but not heap !
				MoveCards(player.Hand, row, player.Hand.FindByName(card.Name))
				MoveCards(player.Deck, row, player.Deck.FindByName(card.Name))
			}
		}
		game.Score()

		if card.Effects.Has(Scorch) {
			if card.Row == Special { //Scorch card
				game.Scorch()
			} else { //Villentremerth: destroys the opponent in its row
				enemy.ScorchAlt(card)
			}
			MoveCard(player.Hand, player.Heap, card)
			enemy.Score(game.WeatherCards.Effects())
		}
		//game.Switch()
	}
	return false
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
	return true
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

func InitDeck(all *CardList, faction Faction) *Cards {
	deck := &Cards{}
	_all := *all
	for id, _ := range _all {
		card := _all[id]
		if card.Faction == faction || card.Faction == Neutral {
			card.Id = id
			card.Faction = faction
			deck.Add(card)
		}
	}
	return deck
}

func (player *GameSide) Score(weather Effects) (sum int) {
	player.ScoreCloseCombat = player.CloseCombat.Score(weather.Has(BitingFrost))
	player.ScoreRangedCombat = player.RangedCombat.Score(weather.Has(ImpenetrableFog))
	player.ScoreSiege = player.Siege.Score(weather.Has(TorrentialRain))
	player.CachedScore = player.ScoreCloseCombat + player.ScoreRangedCombat + player.ScoreSiege
	return player.CachedScore
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
		if card.score > maxScore && !card.Effects.Has(Hero) {
			maxScore = card.score
		}
	}
	for _, card := range cards {
		if card.score == maxScore && !card.Effects.Has(Hero) {
			card.Strength = 0
			card.score = -1
			log.Print("scorched", card)
		}
	}
	game.Clean() // deletes the scorched cards
}

func (gs *GameSide) ScorchAlt(card *Card) {
	rowToScorch := gs.GetRow(card.Row)
	maxScore := 0
	for _, enemy_card := range *rowToScorch {
		if enemy_card.score > maxScore && !card.Effects.Has(Hero) {
			maxScore = enemy_card.score
		}
	}
	for _, enemy_card := range *rowToScorch {
		if enemy_card.score == maxScore && !card.Effects.Has(Hero) {
			MoveCard(rowToScorch, gs.Heap, enemy_card)
		}
	}
}

func (card *Card) EligibleMedic() bool {
	for _, eff := range card.Effects {
		if eff == Hero ||
			eff == Decoy ||
			eff == CommanderHorn ||
			eff == Scorch {
			return false
		}
	}
	return true
}

func (g *Game) RoundWinner() Turn {
	scoreA, scoreB := g.Score()
	if scoreA > scoreB {
		return PlayerA
	}
	if scoreB > scoreA {
		return PlayerB
	}
	return Tie
}

func (g *Game) RoundFinished() bool {
	return g.SideA.Passed && g.SideB.Passed
}

// starts a new round if required, and reports if it started one
func (g *Game) NextRound() bool {
	if g.RoundFinished() {
		round_winner := g.RoundWinner()

		g.History = append(g.History, round_winner)

		if g.Round() < MAX_ROUND {
			g.SideA.EndRound()
			g.SideB.EndRound()
			if round_winner != Tie {
				g.Turn = round_winner //winner starts new round
			}
		}

		return true
	}
	/*if g.Round() > MAX_ROUND {
		panic("max rounds exceeded")
	}*/
	return false
}

func (g *Game) Finished() bool {
	if !g.RoundFinished() {
		return false
	}
	side, max := g.MaxRoundsWon()
	if side == Tie { //game ends in a tie
		return g.Round() == MAX_ROUND
	} else {
		return g.Round() == MAX_ROUND || max > 1
	}
}

func (g *Game) Winner() Turn {
	winner, wins := g.MaxRoundsWon()
	if wins == 2 && winner != Tie {
		return winner
	}
	rounds := g.RoundsWon()

	if rounds[PlayerB] > rounds[PlayerA] {
		return PlayerB
	}
	if rounds[PlayerA] > rounds[PlayerB] {
		return PlayerA
	}
	/*
		if rounds[PlayerA] == rounds[PlayerB]{
				return Tie
		}
	*/
	return Tie
}
