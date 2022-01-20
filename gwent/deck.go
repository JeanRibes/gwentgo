package gwent

type PlayerDeck struct {
	Leader *Leader
	Deck   *CardList
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

type PlayerData struct {
	Pseudo string
	Cookie string
	Decks  *map[Faction]*PlayerDeck
}

func NewPlayerData(pseudo string, cookie string) *PlayerData {
	return &PlayerData{
		Pseudo: pseudo,
		Cookie: cookie,
		Decks:  &map[Faction]*PlayerDeck{},
	}
}
