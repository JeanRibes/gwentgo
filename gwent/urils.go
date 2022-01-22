package gwent

type BondMap map[string]int

func (bond BondMap) Add(card *Card) {
	if _, ok := bond[card.Name]; ok {
		bond[card.Name] += 1
	} else {
		bond[card.Name] = 1
	}
}

func (bond BondMap) Has(card Card) bool {
	_, ok := bond[card.Name]
	return ok
}

func (bond BondMap) Get(card *Card) int {
	return bond[card.Name]
}

func ToCards(list CardList) *CardList {
	out := &CardList{}
	/*for i, card := range list {
		card.Id = i
		out.Add(card)
	}*/
	for _, card := range list {
		out.Add(card)
	}
	return out
}

func CardPointers(in []Card) (out CardList) {
	out = make(CardList, len(in))
	for i, card := range in {
		out[i] = &card
	}
	return out
}

func Creategame() *Game {
	deckP2 := InitDeck(AllCardsList.Copy(), ScoiaTael).Copy()
	handP2 := ToCards(CardList{
		AllCardsMap["saesenthessis"],
		AllCardsMap["isengrim-faoiltiarna"],
		AllCardsMap["filavandrel"],
		AllCardsMap["yaevinn"],
		AllCardsMap["dwarf-skirmisher"],
		AllCardsMap["dwarf-skirmisher"],
		AllCardsMap["havekar-healer"],
		AllCardsMap["scorch"],
		AllCardsMap["biting-frost"],
	}).Copy()
	handP2.SetFaction(ScoiaTael)
	handP2.GetByName("scorch").Strength = 99

	deckP2.Remove(deckP2.GetByName("dwarf-skirmisher"))
	deckP2.Remove(deckP2.GetByName("dwarf-skirmisher"))

	deckP1 := InitDeck(AllCardsList.Copy(), NorthernRealms).Copy()

	handP1 := ToCards(CardList{
		(*AllCardsList)[0],               // #0 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		(*AllCardsList)[1],               // #2 [Northern Realms][CloseCombat] Cirilla Fiona Elen Riannon (15) {Hero,}
		(*AllCardsList)[3],               // #3 [Northern Realms][Weather] Clear Weather (0) {ClearWeather,}
		(*AllCardsList)[9],               // #4 [Northern Realms][CloseCombat] Avallac’h (0) {Spy,Hero,}
		(*AllCardsList)[32],              // #5 [Northern Realms][Siege] Thaler (1) {Spy,}
		(*AllCardsList)[25],              // #7 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		AllCardsMap["villentretenmerth"], // #1 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		AllCardsMap["ves"],               // #6 [Northern Realms][CloseCombat] Poor F*cking Infantry (1) {TightBond,}
		AllCardsMap["scorch"],            // #8 [Northern Realms][CloseCombat] Geralt of Rivia (15) {Hero,}
		(*AllCardsList)[62],              // #9 [Northern Realms][RangedCombat] Yennefer of Vengerberg (7) {Medic,Hero,}
	}).Copy()
	handP1.SetFaction(NorthernRealms)

	//t.Log("deck P2", deckP2)

	return NewGame(deckP1, handP1, deckP2, handP2)
}
