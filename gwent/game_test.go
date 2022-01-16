package gwent

import "testing"

func TestInit(t *testing.T) {
	deck := InitDeck(AllCards, NorthernRealms)
	hand := InitHand(deck)
	t.Log(len(hand))
	for _, card := range hand {
		t.Logf("%s,%v+\n", card.Name, card)
	}
}

func TestList(t *testing.T) {
	hand := ToCards([]Card{
		AllCardsMap["saesenthessis"],
		AllCardsMap["isengrim-faoiltiarna"],
		AllCardsMap["filavandrel"],
		AllCardsMap["yaevinn"],
	}).Reindex()
	list := hand.List()
	for _, lcard := range list {
		hcard, ok := hand[lcard.Id]
		if !ok {
			t.Fatal("map lookup failed")
		}
		if hcard.DisplayName != lcard.DisplayName {
			t.Errorf("%s was %s", hcard, lcard)
			t.Fatal("bad copy")
		}
	}
}

func TestFindByName(t *testing.T) {
	hand := ToCards([]Card{
		AllCardsMap["saesenthessis"],
		AllCardsMap["saesenthessis"],
		AllCardsMap["isengrim-faoiltiarna"],
		AllCardsMap["filavandrel"],
		AllCardsMap["yaevinn"],
	}).Reindex().SetFaction(ScoiaTael)
	hand.FindByName2("yaevinn")
	hand.GetByName("saesenthessis")
	list := hand.FindByName("saesenthessis")
	for _, lcard := range list {
		if lcard.Name != "saesenthessis" {
			t.Fatal("bad name")
		}
		hcard, ok := hand[lcard.Id]
		if !ok {
			t.Fatal("map lookup failed")
		}
		if hcard.DisplayName != lcard.DisplayName {
			t.Errorf("%s was %s", hcard, lcard)
			t.Fatal("bad copy")
		}
	}
}

func TestGame(t *testing.T) {
	deckP1 := InitDeck(AllCards, NorthernRealms).Reindex()

	handP1 := ToCards([]Card{
		AllCards[0],  // #0 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		AllCards[1],  // #2 [Northern Realms][CloseCombat] Cirilla Fiona Elen Riannon (15) {Hero,}
		AllCards[3],  // #3 [Northern Realms][Weather] Clear Weather (0) {ClearWeather,}
		AllCards[9],  // #4 [Northern Realms][CloseCombat] Avallac’h (0) {Spy,Hero,}
		AllCards[32], // #5 [Northern Realms][Siege] Thaler (1) {Spy,}
		AllCards[25], // #7 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		AllCards[45], // #1 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		AllCards[62], // #6 [Northern Realms][CloseCombat] Poor F*cking Infantry (1) {TightBond,}
		AllCards[62], // #8 [Northern Realms][CloseCombat] Geralt of Rivia (15) {Hero,}
		AllCards[62], // #9 [Northern Realms][RangedCombat] Yennefer of Vengerberg (7) {Medic,Hero,}
	}).Reindex().SetFaction(NorthernRealms)

	deckP2 := InitDeck(AllCards, ScoiaTael).Reindex()
	handP2 := ToCards([]Card{
		AllCardsMap["saesenthessis"],
		AllCardsMap["isengrim-faoiltiarna"],
		AllCardsMap["filavandrel"],
		AllCardsMap["yaevinn"],
		AllCardsMap["dwarf-skirmisher"],
		AllCardsMap["dwarf-skirmisher"],
		AllCardsMap["dwarf-skirmisher"],
		AllCardsMap["havekar-healer"],
		AllCardsMap["scorch"],
		AllCardsMap["torrential-rain"],
	}).SetFaction(ScoiaTael).Reindex()

	/*t.Log("hand P2",handP2)
	t.Log("deck P2", deckP2)*/

	game := Game{}.New(deckP1, handP1, deckP2, handP2)
	deckP2.CheckIds()
	handP2.CheckIds()

	first := handP1.GetByName("geralt-of-rivia") /*"yennefer-of-vengerberg"*/
	t.Log("first move:", first)

	game.PlayMove(first, first.Row, &game.SideA, &game.SideB)
	//t.Log(game.String())

	second := handP2.GetByName("dwarf-skirmisher")
	t.Log("second move:", second)
	t.Log(handP2.FindByName2("dwark"))

	game.PlayMove(second, second.Row, &game.SideB, &game.SideA)
	t.Log(game)

}
