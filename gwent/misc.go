package gwent

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type enum uint

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
	case Neutral:
		return "Neutral"
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
	case "Neutral":
		return Neutral, nil
	default:
		println("faction not found", str)
		return 0, errors.New("invalid faction: " + str)
	}
	//return 0, errors.New("eerr ??")
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

func (card *Card) SetFaction(faction Faction) *Card {
	card.Faction = faction
	return card
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

func (card *Card) Copy() *Card {
	_copy := *card
	return &_copy
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
	case Decoy:
		return "Decoy"
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
	case "Decoy":
		return Decoy, nil
	default:
		return 0, errors.New("invalid effect: " + effect)
	}
}
func (effect *Effect) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	*effect, err = EffectFromString(str[1 : len(str)-1])
	return err
}

func MoveCard(source *CardList, dest *CardList, card *Card) {
	dest.Add(card)
	source.Remove(card)
}
func MoveCards(source *CardList, dest *CardList, cards *CardList) {
	source.Removes(cards)
	dest.Adds(cards)
}

func (cards *CardList) Adds(news *CardList) *CardList {
	for _, card := range *news {
		cards.Add(card)
	}
	return cards
}

func (cards *CardList) Removes(to_remove *CardList) *CardList {
	for _, card := range *to_remove {
		cards.Remove(card)
	}
	return cards
}

func (cards *CardList) Effects() Effects {
	effects := Effects{}
	for _, card := range *cards {
		effects = append(effects, card.Effects...)
	}
	return effects
}

func (cards *CardList) FindByName(name string) *CardList {
	ret := CardList{}
	for _, card := range *cards {
		if card.Name == name {
			ret = append(ret, card)
		}
	}
	return &ret
}

func (cards *CardList) FindByName2(name string) (ret CardList) {
	_cards := *cards
	for key, card := range _cards {
		if card.Name == name {
			c := _cards[key]
			ret = append(ret, c)
		}
	}
	return ret
}

func (cards *CardList) GetByName(name string) *Card {
	for _, card := range *cards {
		if card.Name == name {
			return card
		}
	}
	return nil
}

func (cards *CardList) CountByName(name string) (count int) {
	for _, card := range *cards {
		if card.Name == name {
			count += 1
		}
	}
	return count
}

func (cards *CardList) ShallowCopy() *CardList {
	out := CardList{}
	for key, value := range *cards {
		out[key] = value
	}
	return &out
}

func (cards *CardList) SortKeys(start int) int {
	i := start
	for _, card := range *cards {
		card.Id = i
		i += 1
	}
	return i
}

func (cards *CardList) List() *CardList {
	l := make(CardList, len(*cards))
	i := 0
	for _, card := range *cards {
		l[i] = card
		i += 1
	}
	return &l
}

func (cards *CardList) SetFaction(faction Faction) *CardList {
	_cards := *cards
	for key, card := range _cards {
		card.Faction = faction
		_cards[key] = card
	}
	return &_cards
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
func (in *CardList) Copy() *CardList {
	out := make(CardList, len(*in))
	for i, card := range *in {
		out[i] = card.Copy()
	}
	return &out
}
func (list *CardList) String() string {
	s := ""
	for _, card := range *list {
		if card != nil {
			s += card.String() + "\n"
		} else {
			s += "<nil>\n"
		}
	}
	return s
}

func (list *CardList) Cards() *CardList {
	return ToCards(*list)
}

/*
Performs a copy !
*/
func (list *CardList) FilterFaction(faction Faction) *CardList {
	ret := CardList{}
	for _, card := range *list {
		if card.Faction == faction || card.Faction == Neutral {
			ret = append(ret, card.Copy())
		}
	}
	return &ret
}

func (list *CardList) Len() int {
	return len(*list)
}
func (_list *CardList) Less(i, j int) bool {
	list := *_list
	return list[i].Strength < list[j].Strength
}
func (_list *CardList) Swap(i, j int) {
	list := *_list
	list[i], list[j] = list[j], list[i]
}

func (_list *CardList) group() *CardList {
	list := *_list
	maxlen := len(list)
	for index, card := range list {
		if index < maxlen-2 {
			next := list[index+1]
			if len(card.Name) > len(next.Name) && card.Strength == next.Strength {
				list[index] = next
				list[index+1] = card
			}
		}
	}
	return _list
}

func (_list *CardList) IsGroupSorted() bool { // use on a strenght-sorted list
	list := *_list
	strength := -1
	lenstr := -1
	for _, card := range list {
		if card.Strength < strength {
			return false
		} //not sorted by strength

		if card.Strength > strength { //we changed of sorting group
			strength = card.Strength
			lenstr = len(card.Name)
		}

		if lenstr > len(card.Name) {
			return false
		} // not sorted by name length

		lenstr = len(card.Name)
	}
	return true
}

func (list *CardList) GroupSort(max_passes int) *CardList {
	sort.Sort(list)
	n := 0
	for !list.IsGroupSorted() {
		list.group()
		n += 1
		if n > max_passes {
			//log.Printf("gave up sorting after %d passes\n", n)
			return list
		}
	}
	return list
}

func (list *CardList) GetById(id int) *Card {
	for _, card := range *list {
		if card.Id == id {
			return card
		}
	}
	return nil
}

func (list *CardList) Add(card *Card) *CardList {
	*list = append(*list, card)
	return list
}

func (list *CardList) Index(card *Card) int {
	for i, _card := range *list {
		if _card.Id == card.Id {
			return i
		}
	}
	return -1
}

func (list *CardList) Has(card *Card) bool {
	return list.Index(card) >= 0
}

func (list *CardList) Remove(card *Card) *CardList {
	index := list.Index(card)
	if index < 0 {
		panic(fmt.Errorf("index < 0: card %s is not present in the list", card.String()))
	}

	if list.Len() == 1 {
		*list = CardList{}
	} else {
		_list := *list
		*list = append(_list[:index], _list[index+1:]...)
	}
	return list
}

func (list *CardList) CheckNil() bool {
	for _, card := range *list {
		if card == nil {
			return false
		}
	}
	return true
}
func MoveCard_list(source *CardList, dest *CardList, card *Card) {
	if source.Has(card) {
		dest.Add(card)
		source.Remove(card)
	} else {
		panic(fmt.Errorf("illegal move: card #%d not present in source\n", card.Id))
	}
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

func UnstringTurn(str string) Turn {
	if str == "PlayerA" {
		return PlayerA
	}
	if str == "PlayerB" {
		return PlayerB
	}

	return Tie
}

func (turn *Turn) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	str = str[1 : len(str)-1]
	*turn = UnstringTurn(str)
	return nil
}
func (turn Turn) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(turn.String())
	buffer.WriteRune('"')
	return buffer.Bytes(), nil
}

func (card *Card) IsUnit() bool {
	return card.Strength > 0 || card.Effects.Has(Spy) || card.Effects.Has(Medic)
}
