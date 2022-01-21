package gwent

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
	deck := &Cards{}
	faction := playerDeck.Faction()
	all := *playerDeck.Deck
	for _, card := range all {
		if card.Faction == faction || card.Faction == Neutral {
			copied := card.Copy()
			copied.Faction = faction
			deck.Add(copied)
		}
	}
	hand := InitHand(deck)
	return NewGameSide(deck, hand)
}

func (playerDeck *PlayerDeck) AddToDeck(card *Card) {
	MoveCard_list(playerDeck.Rest, playerDeck.Deck, card)
	playerDeck.Deck.GroupSort(1)
}

func (playerDeck *PlayerDeck) RemoveFromDeck(card *Card) {
	MoveCard_list(playerDeck.Deck, playerDeck.Rest, card)
	playerDeck.Rest.GroupSort(1)
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
		OwnedCards: AllCardsList.Copy().Cards().Reindex().List().GroupSort(1),
	}
	playerdata.Decks = &Decks{
		playerdata.NewPlayerDeck(NorthernRealms),
		playerdata.NewPlayerDeck(Nilfgaard),
		playerdata.NewPlayerDeck(ScoiaTael),
		playerdata.NewPlayerDeck(Monsters),
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
