package gwent

import "fmt"

type GameSide struct {
	Hand Cards //what the player sees and can choose to play
	Deck Cards //the rest of the cards, that may be used with spies
	Heap Cards //discarded cards, that may come back to hand with Medics

	CloseCombat Cards

	RangedCombat Cards

	Siege Cards
}

func (g GameSide) New(deck Cards, hand Cards) GameSide {
	return GameSide{
		Deck:         deck,
		Hand:         hand,
		Heap:         Cards{},
		CloseCombat:  Cards{},
		RangedCombat: Cards{},
		Siege:        Cards{},
	}
}

type Game struct {
	WeatherCards Cards
	SideA        GameSide
	SideB        GameSide
}

func (g Game) New(
	deckA Cards,
	handA Cards,
	deckB Cards,
	handB Cards,
) Game {
	g.SideA = GameSide{}.New(deckA, handA)
	g.SideB = GameSide{}.New(deckB, handB)
	g.WeatherCards = Cards{}
	return g
}

func (player *GameSide) GetRow(row Row) Cards {
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

func (player *GameSide) Merge() []*Card {
	arr := make([]*Card,
		len(player.CloseCombat)+len(player.RangedCombat)+len(player.Siege))
	for key, _ := range player.CloseCombat {
		c := player.CloseCombat[key]
		arr = append(arr, &c)
	}
	for key, _ := range player.RangedCombat {
		c := player.RangedCombat[key]
		arr = append(arr, &c)
	}
	for key, _ := range player.Siege {
		c := player.Siege[key]
		arr = append(arr, &c)
	}
	return arr
}

func (game *Game) Merge() []*Card {
	arr := []*Card{}
	arr = append(arr, game.SideA.Merge()...)
	arr = append(arr, game.SideB.Merge()...)
	return arr
}

func (row Cards) Clean() []*Card {
	removed := []*Card{}
	for key, card := range row {
		if card.score < 0 {
			c := row[key]
			removed = append(removed, &c)
			delete(row, card.Id)
		}
	}
	return removed
}

func (row Cards) Clear() {
	for _, card := range row {
		delete(row, card.Id)
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

func (gs GameSide) String() string {
	return fmt.Sprintf("CloseCombat: %s\nRangedCombat: %s\nSiege: %s\n",
		gs.CloseCombat,
		gs.RangedCombat,
		gs.Siege)
}

func (g Game) String() string {
	return fmt.Sprintf("Side A\n%s===========\nWeather: %s\n==========\nSide B\n%s",
		g.SideA, g.WeatherCards, g.SideB)
}
