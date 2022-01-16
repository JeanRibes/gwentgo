package gwent

import "testing"

func TestGame(t *testing.T) {

	deckP2 := InitDeck(AllCards, ScoiaTael).Reindex()
	handP2 := ToCards(CardList{
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
	}.Copy()).SetFaction(ScoiaTael).Reindex()

	deckP1 := InitDeck(AllCards, NorthernRealms).Reindex()

	handP1 := ToCards(CardList{
		&AllCards[0],  // #0 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		&AllCards[1],  // #2 [Northern Realms][CloseCombat] Cirilla Fiona Elen Riannon (15) {Hero,}
		&AllCards[3],  // #3 [Northern Realms][Weather] Clear Weather (0) {ClearWeather,}
		&AllCards[9],  // #4 [Northern Realms][CloseCombat] Avallac’h (0) {Spy,Hero,}
		&AllCards[32], // #5 [Northern Realms][Siege] Thaler (1) {Spy,}
		&AllCards[25], // #7 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		&AllCards[45], // #1 [Northern Realms][CloseCombat] Blue Stripes Commando (4) {TightBond,}
		&AllCards[62], // #6 [Northern Realms][CloseCombat] Poor F*cking Infantry (1) {TightBond,}
		&AllCards[62], // #8 [Northern Realms][CloseCombat] Geralt of Rivia (15) {Hero,}
		&AllCards[62], // #9 [Northern Realms][RangedCombat] Yennefer of Vengerberg (7) {Medic,Hero,}
	}.Copy()).Reindex().SetFaction(NorthernRealms)

	//t.Log("deck P2", deckP2)

	game := NewGame(deckP1, handP1, deckP2, handP2)
	deckP1, handP1, deckP2, handP2 = nil, nil, nil, nil
	/*deckP2.CheckIds()
	handP2.CheckIds()*/

	first := game.SideA.Hand.GetByName("geralt-of-rivia")
	//t.Log("first move:", first)

	game.PlayMove(first, first.Row, game.SideA, game.SideB)
	if game.SideA.CloseCombat.CountByName("geralt-of-rivia") != 1 {
		t.Fail()
	}

	//t.Log(game)
	if card := game.SideA.Hand.GetByName(first.Name); card != nil {
		if card.Id == first.Id {
			t.Fatal("card not removed from hand")
		}
	}

	second := game.SideB.Hand.GetByName("dwarf-skirmisher")
	//t.Log("second move:", second)
	//t.Log(game.SideB.Hand.FindByName2("dwark"))

	game.PlayMove(second, second.Row, game.SideB, game.SideA)
	if game.SideB.CloseCombat.Len() != 3 {
		t.Fatal("Muster didn't work")
	}
	//t.Log(game)

}
