package gwent

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
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

func (faction Faction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(faction.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (faction Faction) String() string {
	switch faction {
	case NorthernRealms:
		return "NorthernRealms"
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
func FactionFromString(str string) (Faction, error) {
	switch str {
	case "NorthernRealms":
		return NorthernRealms, nil
	case "Northern Realms":
		return NorthernRealms, nil
	case "ScoiaTael":
		return ScoiaTael, nil
	case "Nilfgaard":
		return Nilfgaard, nil
	case "Monsters":
		return Monsters, nil
	default:
		println("faction not found", str)
		return 0, errors.New("invalid faction: " + str)
	}
	return 0, errors.New("eerr ??")
}
func (faction *Faction) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	*faction, err = FactionFromString(str[1 : len(str)-1]) //remove quotes
	return err
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

func (row Row) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(row.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func RowFromString(str string) Row {
	switch str {
	case "CloseCombat":
		return CloseCombat
	case "RangedCombat":
		return RangedCombat
	case "Siege":
		return Siege
	case "Weather":
		return Weather
	case "Special":
		return Special
	default:
		return 99
	}
}

func (row *Row) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	*row = RowFromString(str[1 : len(str)-1])
	if *row == 99 {
		return errors.New("invalid Row")
	}
	return nil
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
	Decoy
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
	Image       string
	Strength    int
	Effects     Effects
	score       int //cached score
}

func (card *Card) Picture() string {
	if card.Image != "" {
		return card.Image
	} else {
		picture := strings.Replace(card.DisplayName, " ", "-", -1)
		picture = strings.Replace(picture, "’", "", -1)
		picture = strings.Replace(picture, "*", "u", -1)
		return picture + "-gwent-card.jpg"
	}
}

func (card *Card) BoardRow() string {
	if card.Row != Special {
		return card.Row.String()
	}
	if card.Effects.Has(Scorch) {
		return "Scorch"
	}
	if card.Effects.Has(Agile) {
		return "Agile"
	}
	return "All"
}

func (card *Card) String() string {
	return fmt.Sprintf("#%d [%s][%s] %s (%d->%d) {%s}",
		card.Id,
		card.Faction,
		card.Row,
		card.Name,
		card.Strength,
		card.score,
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
func (eff Effect) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(eff.String())
	buffer.WriteRune('"')
	return buffer.Bytes(), nil
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
func EffectFromString(effect string) (Effect, error) {
	switch effect {
	case "Agile":
		return Agile, nil
	case "CommanderHorn":
		return CommanderHorn, nil
	case "Hero":
		return Hero, nil
	case "Medic":
		return Medic, nil
	case "MoraleBoost":
		return MoraleBoost, nil
	case "Muster":
		return Muster, nil
	case "Spy":
		return Spy, nil
	case "TightBond":
		return TightBond, nil
	case "Scorch":
		return Scorch, nil
	case "ClearWeather":
		return ClearWeather, nil
	case "BitingFrost":
		return BitingFrost, nil
	case "ImpenetrableFog":
		return ImpenetrableFog, nil
	case "TorrentialRain":
		return TorrentialRain, nil
	default:
		return 0, errors.New("invalid effect")
	}
}
func (effect *Effect) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	*effect, err = EffectFromString(str[1 : len(str)-1])
	return err
}

type Cards map[int]*Card

func MoveCard(source *Cards, dest *Cards, card *Card) {
	dest.Add(card)
	source.Remove(card)
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

func (cards *Cards) CheckIds() bool {
	for i, card := range *cards {
		if i != card.Id {
			log.Panicf("ID clash: #%d should be %s", i, card)
			return false
		}
	}
	return true
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

func (cards *Cards) GetById(id int) *Card {
	_cards := *cards
	if card, ok := _cards[id]; ok {
		return card
	} else {
		return nil
	}
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

func (cards *Cards) SortKeys(start int) int {
	backup := cards.Copy()
	cards.Clear()

	i := start
	for _, card := range backup {
		card.Id = i
		i += 1
		cards.Add(card)
	}
	return i
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
func (list CardList) String() string {
	s := ""
	for _, card := range list {
		if card != nil {
			s += card.String() + "\n"
		} else {
			s += "<nil>\n"
		}
	}
	return s
}

type Turn enum

const (
	PlayerA Turn = iota
	PlayerB
	Tie
)

func (turn Turn) String() string {
	if turn == PlayerA {
		return "PlayerA"
	}
	if turn == PlayerB {
		return "PlayerB"
	}
	return "Tie"
}
func (turn *Turn) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	str = str[1 : len(str)-1]
	if str == "PlayerA" {
		*turn = PlayerA
	}
	if str == "PlayerB" {
		*turn = PlayerB
	} else {
		*turn = Tie
	}
	return nil
}
func (turn Turn) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(turn.String())
	buffer.WriteRune('"')
	return buffer.Bytes(), nil
}
