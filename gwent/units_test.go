package gwent

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func genHand() *CardList {
	return ToCards(CardList{
		AllCardsMap["saesenthessis"],
		AllCardsMap["isengrim-faoiltiarna"],
		AllCardsMap["filavandrel"],
		AllCardsMap["yaevinn"],
		AllCardsMap["geralt-of-rivia"],
	}).Copy()
}

func TestCards_List(t *testing.T) {
	hand := genHand()
	list := *hand.List()
	assert.Equal(t, hand.Len(), 5)
	assert.NotEqual(t, list[0].Id, list[1].Id)
}

func TestCards_FindByName(t *testing.T) {
	hand := genHand()
	list := hand.FindByName("saesenthessis")
	assert.Equal(t, (*list)[0].Name, "saesenthessis")
}

func TestCards_GetByName(t *testing.T) {
	hand := genHand()
	name := "yaevinn"
	card := hand.GetByName(name)
	if card == nil {
		t.Fatal("returned nil")
	}
	if card.Name != name {
		t.Fatal("bad name")
	}
}

func TestCards_GetById(t *testing.T) {
	hand := genHand().Copy()
	t.Log(hand)
	card := hand.GetById(106)
	assert.NotNil(t, card)
}

func TestCards_Remove(t *testing.T) {
	hand := genHand()
	geralt := hand.GetByName("geralt-of-rivia")
	len0 := hand.Len()
	hand.Remove(geralt)
	if hand.Len() != len0-1 {
		t.Fail()
	}
}

func TestCards_CountByName(t *testing.T) {
	hand := genHand().Cards()
	nc := hand.GetById(106).Copy()
	hand.Add(nc)
	if hand.CountByName("saesenthessis") != 2 {
		t.Fail()
	}
	if hand.CountByName("geralt-of-rivia") != 1 {
		t.Fail()
	}
}

func TestCards_Has(t *testing.T) {
	hand := genHand()
	geralt := hand.GetByName("geralt-of-rivia")
	if !hand.Has(geralt) {
		t.Fail()
	}
}

func TestCards_Copy(t *testing.T) {
	hand := genHand()
	cop := hand.Copy()
	if cop.Len() < 3 {
		t.Fatal("empty copy")
	}
	geralt := hand.GetByName("geralt-of-rivia")
	hand.Remove(geralt)
	if cop.Len() == hand.Len() {
		t.Fatal("shallow copy")
	}

}

func TestEffects_Has(t *testing.T) {
	card := AllCardsMap["biting-frost"]
	assert.True(t, card.Effects.Has(BitingFrost))

	geralt := AllCardsMap["geralt-of-rivia"]
	assert.True(t, geralt.Effects.Has(Hero))
}

func TestGame_RoundWinner(t *testing.T) {
	game := Game{
		SideA: &GameSide{
			CloseCombat:  &CardList{},
			RangedCombat: &CardList{},
			Siege:        &CardList{},
		},
		SideB: &GameSide{
			CloseCombat:  &CardList{},
			RangedCombat: &CardList{},
			Siege:        &CardList{},
		},
		WeatherCards: &CardList{},
	}
	assert.Equal(t, game.RoundWinner(), Tie)

	game = Game{
		SideA: &GameSide{
			CloseCombat: ToCards(CardList{
				AllCardsMap["ves"].Copy(),
			}),
			RangedCombat: &CardList{},
			Siege:        &CardList{},
		},
		SideB: &GameSide{
			CloseCombat:  ToCards(CardList{AllCardsMap["vreemde"].Copy()}),
			RangedCombat: &CardList{},
			Siege:        &CardList{},
		},
		WeatherCards: &CardList{},
	}
	assert.Equal(t, game.RoundWinner(), PlayerA)
}

func TestGame_RoundFinished(t *testing.T) {
	game := Game{
		SideA:   &GameSide{},
		SideB:   &GameSide{},
		History: []Turn{},
	}
	assert.Equal(t, game.Round(), 0)
	assert.False(t, game.RoundFinished())
	game.SideA.Passed = true
	assert.False(t, game.RoundFinished())
	game.SideB.Passed = true

	assert.True(t, game.RoundFinished())
}

func TestGame_NextRound(t *testing.T) {
	game := NewGame(&CardList{}, &CardList{}, &CardList{}, &CardList{})
	game.SideA.Passed = true
	game.SideB.Passed = true
	game.NextRound()
	assert.Equal(t, game.Round(), 1)
	assert.Equal(t, game.History[0], Tie)
}

func TestGame_Finished(t *testing.T) {
	game := NewGame(&CardList{}, &CardList{}, &CardList{}, &CardList{})
	assert.False(t, game.RoundFinished())
	game.SideA.Passed = true
	game.SideB.Passed = true
	game.History = []Turn{Tie, Tie}
	assert.True(t, game.RoundFinished())
	assert.False(t, game.Finished())

	game.History = []Turn{PlayerA, PlayerA, Tie}
	t.Log(game.MaxRoundsWon())
	t.Log(game.RoundFinished())
	assert.True(t, game.Finished())
	assert.Equal(t, PlayerA, game.Winner())

	game.History = []Turn{PlayerA, PlayerB, Tie}
	assert.True(t, game.Finished())
	assert.Equal(t, Tie, game.Winner())

	game.History = []Turn{PlayerA, PlayerB, PlayerA}
	assert.True(t, game.Finished())
	assert.Equal(t, PlayerA, game.Winner())
}

func TestGame_Winner(t *testing.T) {
	game := Game{
		History: []Turn{PlayerA, PlayerA, Tie},
	}
	t.Log(game.RoundsWon())
	assert.Equal(t, game.Winner(), PlayerA)
}

func TestCard_Copy(t *testing.T) {
	card := AllCardsMap["geralt-of-rivia"]
	ccopy := card.Copy()
	assert.Equal(t, card.String(), ccopy.String())
	assert.Equal(t, card.Strength, ccopy.Strength)
	ccopy.Strength = -2
	assert.NotEqual(t, card.Strength, ccopy.Strength)

}

func TestCardList_GroupSort(t *testing.T) {
	list := &CardList{
		AllCardsMap["saesenthessis"].Copy(),
		AllCardsMap["villentretenmerth"].Copy(),
		AllCardsMap["isengrim-faoiltiarna"].Copy(),
		AllCardsMap["filavandrel"].Copy(),
		AllCardsMap["filavandrel"].Copy(),
		AllCardsMap["yaevinn"].Copy(),
		AllCardsMap["yaevinn"].Copy(),
		AllCardsMap["dwarf-skirmisher"].Copy(),
		AllCardsMap["dwarf-skirmisher"].Copy(),
		AllCardsMap["biting-frost"].Copy(),
		AllCardsMap["ves"].Copy(),
	}
	assert.False(t, sort.IsSorted(list))
	bak := list.GroupSort(20)
	assert.True(t, sort.IsSorted(bak))
	assert.True(t, sort.IsSorted(list))
	assert.True(t, list.IsGroupSorted())
	assert.True(t, bak.IsGroupSorted())
	assert.Equal(t, list.String(), bak.String())
}

func TestPlayerDeck_AddToDeck(t *testing.T) {
	deck := NewPlayerData("", "").Decks.GetByIndex(0)
	assert.Equal(t, deck.Faction(), NorthernRealms)
	_rest := *deck.Rest
	card := _rest[50]
	t.Log(card)

	rlen0 := deck.Rest.Len()
	dlen0 := deck.Deck.Len()

	deck.AddToDeck(card)

	rlen1 := deck.Rest.Len()
	dlen1 := deck.Deck.Len()

	assert.Equal(t, rlen0, rlen1+1)
	assert.Equal(t, dlen0+1, dlen1)
	assert.True(t, deck.Rest.CheckNil())
	assert.True(t, deck.Deck.CheckNil())

	deck.RemoveFromDeck(card)

	assert.True(t, deck.Rest.CheckNil())
	assert.True(t, deck.Deck.CheckNil())
}

func TestCardList_SortKeys(t *testing.T) {
	hand := genHand()
	//t.Log(hand)
	hand.SortKeys(10)
	//t.Log(hand)
	for _, card := range *hand {
		assert.Greaterf(t, card.Id, 9, "sorting ids failed")
	}
}

func TestGame_LivesLeft(t *testing.T) {
	game := Game{
		SideA:   &GameSide{},
		SideB:   &GameSide{},
		History: []Turn{},
	}
	assert.Equal(t, 2, game.LivesLeft(PlayerA))
	game.History = []Turn{Tie}
	assert.Equal(t, 1, game.LivesLeft(PlayerA))
	assert.Equal(t, 1, game.LivesLeft(PlayerB))
	game.History = []Turn{Tie, PlayerA}
	assert.Equal(t, 0, game.LivesLeft(PlayerB))
	assert.Equal(t, 1, game.LivesLeft(PlayerA))

	game.History = []Turn{PlayerA, PlayerA}
	assert.Equal(t, 0, game.LivesLeft(PlayerB))
	assert.Equal(t, 2, game.LivesLeft(PlayerA))

}
