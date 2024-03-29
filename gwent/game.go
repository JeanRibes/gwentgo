package gwent

import (
	"encoding/json"
	"fmt"
	"log"
)

type GameSide struct {
	Leader *Leader
	Hand   *CardList //what the player sees and can choose to play
	Deck   *CardList //the rest of the cards, that may be used with spies
	Heap   *CardList //discarded cards, that may come back to hand with Medics

	CloseCombat  *CardList
	RangedCombat *CardList
	Siege        *CardList

	ScoreCloseCombat  int
	ScoreRangedCombat int
	ScoreSiege        int
	CachedScore       int

	Passed      bool
	MedicAction bool
	Side        Turn
}

func NewGameSide(deck *CardList, hand *CardList) *GameSide {
	return &GameSide{
		Deck:         deck,
		Hand:         hand,
		Heap:         &CardList{},
		CloseCombat:  &CardList{},
		RangedCombat: &CardList{},
		Siege:        &CardList{},
		Passed:       false,
		Side:         Tie,
		Leader:       Leaders["foltest-king-of-temeria"],
	}
}

type Game struct {
	WeatherCards *CardList
	SideA        *GameSide
	SideB        *GameSide

	Turn      Turn
	History   []Turn
	Recording []Record `json:"-"`
}

func NewGame(
	deckA *CardList,
	handA *CardList,
	deckB *CardList,
	handB *CardList,
) *Game {
	sa := NewGameSide(deckA, handA)
	sb := NewGameSide(deckB, handB)
	sa.Side = PlayerA
	sb.Side = PlayerB
	return &Game{
		WeatherCards: &CardList{},
		SideA:        sa,
		SideB:        sb,
		History:      []Turn{},
		Turn:         PlayerA,
		Recording:    []Record{},
	}
}

func GameFromDecks(a *PlayerDeck, b *PlayerDeck) (*Game, error) {
	if !a.Eligible() {
		return nil, fmt.Errorf("deck %s is ineligible for gameplay", a.Name)
	}
	if !b.Eligible() {
		return nil, fmt.Errorf("deck %s is ineligible for gameplay", b.Name)
	}
	ah, ad := a.DrawHandDeck()
	bh, bd := b.DrawHandDeck()
	game := NewGame(ad, ah, bd, bh)
	game.SideA.Leader = a.Leader
	game.SideB.Leader = b.Leader
	game.Sort()
	return game, nil
}

func (g *Game) fixSides() {
	g.SideA.Side = PlayerA
	g.SideB.Side = PlayerB
}

func (g *Game) Pass(side *GameSide) {
	//TODO: use a Turn instead of a GameSide
	g.recordPass(side.Side)
	side.Passed = true
	if !g.NextRound() { //if only one player passed, switch turn
		g.Switch()
	}
}

func (g *Game) Round() int {
	return len(g.History)
}

func (player *GameSide) GetRow(row Row) *CardList {
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

func (game *Game) Merge() *CardList {
	arr := CardList{}
	arr = append(arr, game.SideA.Merge()...)
	arr = append(arr, game.SideB.Merge()...)
	return &arr
}

func (row *CardList) Clean() *CardList {
	out := CardList{}
	row.Clear()
	for _, card := range *row.Copy() {
		if card.score == -1 {
			out.Add(card)
		} else {
			row.Add(card)
		}
	}
	return &out
}

func (row *CardList) Clear() {
	*row = CardList{}
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
	return g.Side(g.Turn)
}
func (g *Game) Enemy() *GameSide {
	return g.Side(g.SideEnemy())
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
func (g *Game) Side(side Turn) *GameSide {
	if side == PlayerA {
		return g.SideA
	}
	if side == PlayerB {
		return g.SideB
	}
	return nil
}
func (g *Game) EnemySide(side Turn) Turn {
	if side == PlayerA {
		return PlayerB
	}
	if side == PlayerB {
		return PlayerA
	}
	return Tie
}
func (g *Game) Switch() *Game {
	g.fixSides()
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

func (g *Game) LivesLeft(side Turn) int {
	lives := 2
	for _, who_won := range g.History {
		if who_won != side {
			lives -= 1
		}
	}
	return lives
}

func (g *Game) DeepCopy() *Game {
	data, err := json.Marshal(g)
	if err != nil {
		log.Println(err)
		return nil
	}
	var copied Game
	if err := json.Unmarshal(data, &copied); err != nil {
		log.Println("unable to unmarshal", err)
		log.Printf("%s", data)
		return nil
	}
	//copied.Turn=g.Turn
	return &copied
}
