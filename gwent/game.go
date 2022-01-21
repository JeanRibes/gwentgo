package gwent

import (
	"fmt"
)

type GameSide struct {
	Hand *Cards //what the player sees and can choose to play
	Deck *Cards //the rest of the cards, that may be used with spies
	Heap *Cards //discarded cards, that may come back to hand with Medics

	CloseCombat  *Cards
	RangedCombat *Cards
	Siege        *Cards

	ScoreCloseCombat  int
	ScoreRangedCombat int
	ScoreSiege        int
	CachedScore       int

	Passed bool
}

func NewGameSide(deck *Cards, hand *Cards) *GameSide {
	return &GameSide{
		Deck:         deck,
		Hand:         hand,
		Heap:         &Cards{},
		CloseCombat:  &Cards{},
		RangedCombat: &Cards{},
		Siege:        &Cards{},
		Passed:       false,
	}
}

type Game struct {
	WeatherCards *Cards
	SideA        *GameSide
	SideB        *GameSide

	Turn    Turn
	History []Turn
}

func NewGame(
	deckA *Cards,
	handA *Cards,
	deckB *Cards,
	handB *Cards,
) *Game {
	return &Game{
		WeatherCards: &Cards{},
		SideA:        NewGameSide(deckA, handA),
		SideB:        NewGameSide(deckB, handB),
		History:      []Turn{},
		Turn:         PlayerA,
	}
}

func (g *Game) Pass(side *GameSide) {
	side.Passed = true
	if !g.NextRound() { //if only one player passed, switch turn
		g.Switch()
	}
}

func (g *Game) Round() int {
	return len(g.History)
}

func (player *GameSide) GetRow(row Row) *Cards {
	switch row {
	case CloseCombat:
		return player.CloseCombat
	case RangedCombat:
		return player.RangedCombat
	case Siege:
		return player.Siege
	default:
		return nil
	}
}

func (player *GameSide) Merge() CardList {
	arr := CardList{}
	for _, card := range *player.CloseCombat {
		arr = append(arr, card)
	}
	for _, card := range *player.RangedCombat {
		arr = append(arr, card)
	}
	for _, card := range *player.Siege {
		arr = append(arr, card)
	}
	return arr
}

func (game *Game) Merge() CardList {
	arr := CardList{}
	arr = append(arr, game.SideA.Merge()...)
	arr = append(arr, game.SideB.Merge()...)
	return arr
}

func (row *Cards) Clean() *CardList {
	removed := CardList{}
	for _, card := range *row {
		if card.score < 0 {
			removed = append(removed, card)
			delete(*row, card.Id)
		}
	}
	return &removed
}

func (row *Cards) Clear() {
	for _, card := range *row {
		delete(*row, card.Id)
	}
}

func (player *GameSide) Clean() {
	player.Heap.Adds(player.CloseCombat.Clean())
	player.Heap.Adds(player.RangedCombat.Clean())
	player.Heap.Adds(player.Siege.Clean())
}
func (game *Game) Clean() {
	game.SideB.Clean()
	game.SideA.Clean()
}

func (game *Game) WeatherClean() {
	weather := game.WeatherCards.Effects()
	if weather.Has(ClearWeather) {
		game.WeatherCards.Clear()
	}
}

func (gs *GameSide) String() string {
	return fmt.Sprintf("CloseCombat: %s\nRangedCombat: %s\nSiege: %s\n",
		gs.CloseCombat,
		gs.RangedCombat,
		gs.Siege)
}

func (g *Game) String() string {
	return fmt.Sprintf("Side A\n%s===========\nWeather: %s\n==========\nSide B\n%s",
		g.SideA, g.WeatherCards, g.SideB)
}

func (g *Game) Check() bool {
	return g.SideA.Check() && g.SideB.Check()
}
func (gs *GameSide) Check() bool {
	return gs.Heap.CheckIds() &&
		gs.Hand.CheckIds() &&
		gs.Deck.CheckIds() &&
		gs.Siege.CheckIds() &&
		gs.RangedCombat.CheckIds() &&
		gs.CloseCombat.CheckIds()
}

func (g *Game) Score() (a int, b int) {
	eff := g.WeatherCards.Effects()
	return g.SideA.Score(eff), g.SideB.Score(eff)
}
func (g *Game) ScoreA() int {
	return g.SideA.Score(g.WeatherCards.Effects())
}
func (g *Game) ScoreB() int {
	return g.SideB.Score(g.WeatherCards.Effects())
}

/*
Ensures that all the IDs in that game are unique, so there are no clashed with Spies and weather cards
*/
func (g *Game) Sort() int {
	index := g.SideA.Hand.SortKeys(0)
	index = g.SideA.Deck.SortKeys(index)
	index = g.SideA.CloseCombat.SortKeys(index)
	index = g.SideA.RangedCombat.SortKeys(index)
	index = g.SideA.Siege.SortKeys(index)

	index = g.WeatherCards.SortKeys(index)

	index = g.SideB.Hand.SortKeys(index)
	index = g.SideB.Deck.SortKeys(index)
	index = g.SideB.CloseCombat.SortKeys(index)
	index = g.SideB.RangedCombat.SortKeys(index)
	index = g.SideB.Siege.SortKeys(index)
	return index
}

func (g *Game) GetCardById(id int) *Card {
	if card := g.SideA.GetCardById(id); card != nil {
		return card
	}
	if card := g.SideB.GetCardById(id); card != nil {
		return card
	}
	return nil
}

func (gs *GameSide) GetCardById(id int) *Card {
	if card := gs.Hand.GetById(id); card != nil {
		return card
	}
	if card := gs.Heap.GetById(id); card != nil {
		return card
	}
	if card := gs.CloseCombat.GetById(id); card != nil {
		return card
	}
	if card := gs.RangedCombat.GetById(id); card != nil {
		return card
	}
	if card := gs.Siege.GetById(id); card != nil {
		return card
	}
	return nil
}

func (g *Game) Player() *GameSide {
	if g.Turn == PlayerA {
		return g.SideA
	}
	if g.Turn == PlayerB {
		return g.SideB
	}
	return nil
}
func (g *Game) Enemy() *GameSide {
	if g.Turn == PlayerA {
		return g.SideB
	}
	if g.Turn == PlayerB {
		return g.SideA
	}
	return nil
}
func (g *Game) SideEnemy() Turn {
	if g.Turn == PlayerA {
		return PlayerB
	}
	if g.Turn == PlayerB {
		return PlayerA
	}
	return Tie
}
func (g *Game) Switch() *Game {
	if g.SideA.Passed {
		g.Turn = PlayerB
		return g
	}
	if g.SideB.Passed {
		g.Turn = PlayerA
		return g
	}
	if g.Turn == PlayerA {
		g.Turn = PlayerB
		return g
	}
	if g.Turn == PlayerB {
		g.Turn = PlayerA
		return g
	}
	return g
}

func (gs *GameSide) EndRound() {
	MoveCards(gs.CloseCombat, gs.Heap, gs.CloseCombat.List())
	MoveCards(gs.RangedCombat, gs.Heap, gs.RangedCombat.List())
	MoveCards(gs.Siege, gs.Heap, gs.Siege.List())
	gs.Passed = false
}

func (g *Game) RoundsWon() map[Turn]int {
	wins := map[Turn]int{
		PlayerA: 0,
		PlayerB: 0,
		Tie:     0,
	}
	for _, winner_round := range g.History { //count rounds won
		wins[winner_round] += 1
	}
	return wins
}
func (g *Game) RoundsWonBy(side Turn) int {
	return g.RoundsWon()[side]
}
func (g *Game) MaxRoundsWon() (_ Turn, rounds_won int) {
	wons := g.RoundsWon()
	rounds_won = 0
	for _, won := range wons {
		if won > rounds_won {
			rounds_won = won
		}
	}
	for side, won := range wons {
		if won == rounds_won {
			return side, rounds_won
		}
	}
	return Tie, -1
}

func (g *Game) LivesLeft(side Turn) (lives int) {
	lives = 2
	for _, who_won := range g.History {
		if who_won != side {
			lives -= 2
		}
	}
	return lives
}
