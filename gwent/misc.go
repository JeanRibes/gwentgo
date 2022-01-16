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

func (card Card) String() string {
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

type Cards map[int]Card

func MoveCard(source Cards, dest Cards, card *Card) {
	source.Remove(card)
	dest.Add(*card)

	source.CheckIds()
	dest.CheckIds()
}
func MoveCards(source Cards, dest Cards, cards []*Card) {
	source.Removes(cards)
	dest.Adds(cards)
}

func (cards Cards) Add(card Card) {
	cards.CheckIds()
	for i, _card := range cards {
		if i == card.Id {
			log.Panicf("ID error: %s clashes with %s", card, _card)
		}
	}

	cards[card.Id] = card
}

func (cards Cards) CheckIds() {
	for i, card := range cards {
		if i != card.Id {
			log.Panicf("ID clash: #%d should be %s", i, card)
		}
	}
}

func (cards Cards) Adds(news []*Card) {
	for _, card := range news {
		cards.Add(*card)
	}
}
func (cards Cards) Remove(card *Card) {
	delete(cards, card.Id)
}

func (cards Cards) Removes(to_remove []*Card) {
	for _, card := range to_remove {
		cards.Remove(card)
	}
}

func (cards Cards) Effects() Effects {
	effects := Effects{}
	for _, card := range cards {
		effects = append(effects, card.Effects...)
	}
	return effects
}

func (cards Cards) FindByName(name string) (ret []*Card) {
	ret = []*Card{}
	for _, card := range cards {
		if card.Name == name {
			ret = append(ret, &card)
		}
		/*if card.Name == name {
			a:=&card
			log.Println(a)
			ret = append(ret, a)
		}*/
	}
	return ret
}

func (cards Cards) FindByName2(name string) (ret []Card) {
	for key, card := range cards {
		if card.Name == name {
			c := cards[key]
			ret = append(ret, c)
		}
	}
	return ret
}

func (cards Cards) GetByName(name string) *Card {
	for key, card := range cards {
		if card.Name == name {
			c := cards[key]
			return &c
		}
	}
	return nil
}

func (cards Cards) Has(card *Card) bool {
	_, ok := cards[card.Id]
	return ok
}

func (cards Cards) Reindex() Cards {
	i := 0
	out := Cards{}
	for key, card := range cards {
		card.Id = i
		c := cards[key]
		out.Add(c)
		i += 1
	}
	return out
}

func (cards Cards) String() string {
	s := ""
	for _, card := range cards {
		s += card.String() + "\n"
	}
	return s
}

func (cards Cards) List() []*Card {
	l := make([]*Card, len(cards))
	i := 0
	for index, _ := range cards {
		c := cards[index]
		l[i] = &c
		i += 1
	}
	return l
}

func (cards Cards) SetFaction(faction Faction) Cards {
	for i, card := range cards {
		card.Faction = faction
		cards[i] = card
	}
	return cards
}

/*type CardSlot enum

const (
	CloseCombatRow CardSlot = iota
	RangedCombatRow
	SiegeRow
	WeatherRow
)*/
