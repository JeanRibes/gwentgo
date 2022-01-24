package gwent

import "fmt"

type PlayerDeck struct {
	Name   string
	Leader *Leader
	Deck   *CardList
	Rest   *CardList
}

func (playerDeck *PlayerDeck) Faction() Faction {
	return playerDeck.Leader.Faction
}

func (playerDeck *PlayerDeck) ToGameSide() *GameSide {
	deck := &CardList{}
	faction := playerDeck.Faction()
	all := *playerDeck.Deck
	for _, card := range all {
		if card.Faction == faction || card.Faction == Neutral {
			deck.Add(card.Copy().SetFaction(faction))
		}
	}
	hand := InitHand(deck)
	return NewGameSide(deck, hand)
}

func (playerDeck *PlayerDeck) AddToDeck(card *Card) {
	MoveCard_list(playerDeck.Rest, playerDeck.Deck, card)
}

func (playerDeck *PlayerDeck) RemoveFromDeck(card *Card) {
	MoveCard_list(playerDeck.Deck, playerDeck.Rest, card)
}

type Decks []*PlayerDeck

type PlayerData struct {
	Pseudo     string
	Cookie     string
	Decks      *Decks
	OwnedCards *CardList
}

func (player *PlayerData) NewPlayerDeck(faction Faction) *PlayerDeck {
	var leader Leader
	if faction == NorthernRealms {
		leader = *Leaders["foltest-king-of-temeria"]
	}
	if faction == Nilfgaard {
		leader = *Leaders["emhyr-var-emreis-his-imperial-majesty"]
	}
	if faction == ScoiaTael {
		leader = *Leaders["francesca-findabair-pureblood-elf"]
	}
	if faction == Monsters {
		leader = *Leaders["eredin-commander-of-the-red-riders"]
	}
	return &PlayerDeck{
		Leader: &leader, // copie ?
		Deck:   &CardList{},
		Name:   faction.String(),
		Rest:   player.OwnedCards.FilterFaction(faction).GroupSort(30),
	}
}

func NewPlayerData(pseudo string, cookie string) *PlayerData {
	playerdata := &PlayerData{
		Pseudo:     pseudo,
		Cookie:     cookie,
		OwnedCards: AllCardsList.Copy().GroupSort(1),
	}
	playerdata.Decks = &Decks{
		playerdata.NewPlayerDeck(NorthernRealms),
		/*playerdata.NewPlayerDeck(Nilfgaard),
		playerdata.NewPlayerDeck(ScoiaTael),
		playerdata.NewPlayerDeck(Monsters),*/
	}
	return playerdata
}

func (decks *Decks) GetByName(name string) *PlayerDeck {
	for _, deck := range *decks {
		if deck.Name == name {
			return deck
		}
	}
	return nil
}

func (_decks *Decks) GetByIndex(index int) *PlayerDeck {
	decks := *_decks
	if index > len(decks)-1 {
		return nil
	}
	return decks[index]
}

func (deck *PlayerDeck) Fill() *PlayerDeck {
	/*for _, card := range *deck.Rest {
		log.Println(card.String())
		deck.AddToDeck(card)
	}*/
	rest := *deck.Rest.Copy()
	for i := 0; i < len(rest); i++ {
		deck.AddToDeck(rest[i])
	}
	return deck
}

func (playerDeck *PlayerDeck) DrawHandDeck() (*CardList, *CardList) {
	n := 10
	if playerDeck.Deck.Len() < n {
		panic(fmt.Errorf("deck %s is too small (%d)", playerDeck.Name, playerDeck.Deck.Len()))
	}
	hand := CardList{}
	deck := playerDeck.Deck.Copy()

	for hand.Len() < n {
		hand.Draw(deck)
	}

	return hand.GroupSort(10), deck.GroupSort(10)
}

func (playerDeck *PlayerDeck) Eligible() bool {
	if playerDeck.Deck.Len() < 22 {
		return false
	}
	unit_cards := 0
	for _, card := range *playerDeck.Deck {
		if card.IsUnit() {
			//only Unit card have a strength > 0 (except Avallach and Havekar Healers)
			unit_cards += 1
		}
	}
	return unit_cards >= 22
}
