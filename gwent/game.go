package gwent

import "fmt"

type GameSide struct {
	Hand *Cards //what the player sees and can choose to play
	Deck *Cards //the rest of the cards, that may be used with spies
	Heap *Cards //discarded cards, that may come back to hand with Medics

	CloseCombat *Cards

	RangedCombat *Cards

	Siege *Cards
}

func NewGameSide(deck *Cards, hand *Cards) *GameSide {
	return &GameSide{
		Deck:         deck,
		Hand:         hand,
		Heap:         &Cards{},
		CloseCombat:  &Cards{},
		RangedCombat: &Cards{},
		Siege:        &Cards{},
	}
}

type Game struct {
	WeatherCards *Cards
	SideA        *GameSide
	SideB        *GameSide
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
	}
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
	arr := make(CardList,
		len(*player.CloseCombat)+len(*player.RangedCombat)+len(*player.Siege))
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

func (row *Cards) Clean() CardList {
	removed := CardList{}
	for _, card := range *row {
		if card.score < 0 {
			removed = append(removed, card)
			delete(*row, card.Id)
		}
	}
	return removed
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
