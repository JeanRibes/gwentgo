package gwent

import (
	"fmt"
	"log"
)

type enum uint

const MAX_CARDS_IN_ROW = 20

type Faction enum

const (
	NorthernRealms Faction = iota
	ScoiaTael
	Nilfgaard
	Monsters
	Neutral
)

func (faction Faction) String() string {
	switch faction {
	case NorthernRealms:
		return "Northern Realms"
	case ScoiaTael:
		return "ScoiaTael"
	case Nilfgaard:
		return "Nilfgaard"
	case Monsters:
		return "Monsters"
	default:
		return "Default"
	}
}

type Row enum

const (
	CloseCombat Row = iota
	RangedCombat
	Siege
	Weather
	Special
)

func (row Row) String() string {
	switch row {

	case CloseCombat:
		return "CloseCombat"
	case RangedCombat:
		return "RangedCombat"
	case Siege:
		return "Siege"
	case Weather:
		return "Weather"
	case Special:
		return "Special"
	default:
		return "error!!"
	}
}

type Effect enum

const (
	Agile         Effect = iota
	CommanderHorn        // double row score (except heroes)
	Hero
	Medic       // draw 1 from heap to hand, no heroes but you can choose !
	MoraleBoost //+1 to cards in row
	Muster
	Spy
	TightBond
	// neutral cards
	/*Decoy*/
	Scorch
	// weather
	ClearWeather
	BitingFrost     //all Close Combat cards = 1 (except Heroes)
	ImpenetrableFog //all Ranged Combat cards = 1 (except Heroes)
	TorrentialRain  //all Siege cards = 1 (except Heroes)
)

type Effects []Effect

func (effects Effects) Has(effect Effect) bool {
	for _, _e := range effects {
		if _e == effect {
			return true
		}
	}
	return false
}

func (effects Effects) Clean() Effects {
	m := map[Effect]bool{}
	for _, effect := range effects {
		m[effect] = true
	}
	ne := Effects{}
	for effect, _ := range m {
		ne = append(ne, effect)
	}
	return ne
}

type Card struct {
	Id          int //to differentiate identical cards
	Faction     Faction
	Row         Row
	Name        string
	DisplayName string
	Strengh     int
	Effects     Effects
	score       int //cached score
}

func (card *Card) String() string {
	return fmt.Sprintf("#%d [%s][%s] %s (%d) {%s}",
		card.Id,
		card.Faction,
		card.Row,
		card.Name,
		card.Strengh,
		card.Effects)
}

func (eff Effect) String() string {

	switch eff {
	case Agile:
		return "Agile"
	case CommanderHorn:
		return "CommanderHorn"
	case Hero:
		return "Hero"
	case Medic:
		return "Medic"
	case MoraleBoost:
		return "MoraleBoost"
	case Muster:
		return "Muster"
	case Spy:
		return "Spy"
	case TightBond:
		return "TightBond"
	case Scorch:
		return "Scorch"
	case ClearWeather:
		return "ClearWeather"
	case BitingFrost:
		return "BitingFrost"
	case ImpenetrableFog:
		return "ImpenetrableFog"
	case TorrentialRain:
		return "TorrentialRain"
	default:
		return "error!"
	}
}

func (eff Effects) String() string {
	s := ""
	if len(eff) == 0 {
		return ""
	}
	for _, effect := range eff {
		s += effect.String() + ","
	}
	return s
}

type Cards map[int]*Card

func MoveCard(source *Cards, dest *Cards, card *Card) {
	source.Remove(card)
	dest.Add(card)
}
func MoveCards(source *Cards, dest *Cards, cards CardList) {
	source.Removes(cards)
	dest.Adds(cards)
}

func (cards *Cards) Add(card *Card) *Cards {
	/*for i, _card := range *cards {
		if i == card.Id {
			log.Panicf("ID error: %s clashes with %s", card, _card)
		}
	}*/
	_cards := *cards
	_cards[card.Id] = card
	return cards
}

func (cards *Cards) CheckIds() *Cards {
	for i, card := range *cards {
		if i != card.Id {
			log.Printf("ID clash: #%d should be %s", i, card)
		}
	}
	return cards
}

func (cards *Cards) Adds(news CardList) *Cards {
	for _, card := range news {
		cards.Add(card)
	}
	return cards
}
func (cards *Cards) Remove(card *Card) *Cards {
	delete(*cards, card.Id)
	return cards
}

func (cards *Cards) Removes(to_remove CardList) *Cards {
	for _, card := range to_remove {
		cards.Remove(card)
	}
	return cards
}

func (cards *Cards) Effects() Effects {
	effects := Effects{}
	for _, card := range *cards {
		effects = append(effects, card.Effects...)
	}
	return effects
}

func (cards *Cards) FindByName(name string) (ret CardList) {
	ret = CardList{}
	for _, card := range *cards {
		if card.Name == name {
			ret = append(ret, card)
		}
	}
	return ret
}

func (cards *Cards) FindByName2(name string) (ret CardList) {
	_cards := *cards
	for key, card := range _cards {
		if card.Name == name {
			c := _cards[key]
			ret = append(ret, c)
		}
	}
	return ret
}

func (cards *Cards) GetByName(name string) *Card {
	for _, card := range *cards {
		if card.Name == name {
			return card
		}
	}
	return nil
}

func (cards *Cards) CountByName(name string) (count int) {
	for _, card := range *cards {
		if card.Name == name {
			count += 1
		}
	}
	return count
}

func (cards *Cards) Has(card *Card) bool {
	_cards := *cards
	_, ok := _cards[card.Id]
	return ok
}

func (cards *Cards) Copy() Cards {
	out := Cards{}
	for key, value := range *cards {
		out[key] = value
	}
	return out
}

func (cards *Cards) Reindex() *Cards {
	backup := cards.Copy()
	cards.Clear()

	i := 0
	for _, card := range backup {
		card.Id = i
		i += 1
		cards.Add(card)
	}
	return cards
}

func (cards *Cards) String() string {
	s := ""
	for _, card := range *cards {
		s += card.String() + "\n"
	}
	return s
}

func (cards *Cards) List() CardList {
	l := make(CardList, len(*cards))
	i := 0
	for _, card := range *cards {
		l[i] = card
		i += 1
	}
	return l
}

func (cards *Cards) SetFaction(faction Faction) *Cards {
	_cards := *cards
	for key, card := range _cards {
		card.Faction = faction
		_cards[key] = card
	}
	return &_cards
}

func (cards *Cards) Len() int {
	return len(*cards)
}

/*type CardSlot enum

const (
	CloseCombatRow CardSlot = iota
	RangedCombatRow
	SiegeRow
	WeatherRow
)*/

type CardList []*Card

/*
Deep Copy
*/
func (in CardList) Copy() CardList {
	out := make(CardList, len(in))
	for i, cp := range in {
		card := *cp
		out[i] = &card
	}
	return out
}
