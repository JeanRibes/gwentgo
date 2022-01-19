package gwent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func genHand() *Cards {
	return ToCards(CardList{
		&AllCards[0],
		&AllCards[62],
		AllCardsMap["saesenthessis"],
		AllCardsMap["saesenthessis"],
		AllCardsMap["isengrim-faoiltiarna"],
		AllCardsMap["filavandrel"],
		AllCardsMap["yaevinn"],
	}.Copy()).Reindex()
}

func TestCards_List(t *testing.T) {
	hand := genHand()
	list := hand.List()
	if hand.Len() != 7 || list[0].Id == list[1].Id {
		t.Fail()
	}

	_hand := *hand
	for _, lcard := range list {
		hcard, ok := _hand[lcard.Id]
		if !ok {
			t.Fatal("map lookup failed")
		}
		if hcard.DisplayName != lcard.DisplayName {
			t.Errorf("%s was %s", hcard, lcard)
			t.Fatal("bad copy")
		}
	}
}

func TestCards_FindByName(t *testing.T) {
	hand := genHand()
	list := hand.FindByName("saesenthessis")
	_hand := *hand
	for _, lcard := range list {
		if lcard.Name != "saesenthessis" {
			t.Fatal("bad name")
		}
		hcard, ok := _hand[lcard.Id]
		if !ok {
			t.Fatal("map lookup failed")
		}
		if hcard.DisplayName != lcard.DisplayName {
			t.Errorf("%s was %s", hcard, lcard)
			t.Fatal("bad copy")
		}
	}
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

func TestCards_Remove(t *testing.T) {
	hand := genHand()
	geralt := hand.GetByName("geralt-of-rivia")
	hand.Remove(geralt)
	if hand.Len() != 6 {
		t.Fail()
	}
}

func TestCards_CountByName(t *testing.T) {
	hand := genHand()
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
}

func TestGame_RoundWinner(t *testing.T) {
	game := Game{
		SideA: &GameSide{
			CloseCombat:  &Cards{},
			RangedCombat: &Cards{},
			Siege:        &Cards{},
		},
		SideB: &GameSide{
			CloseCombat:  &Cards{},
			RangedCombat: &Cards{},
			Siege:        &Cards{},
		},
		WeatherCards: &Cards{},
	}
	assert.Equal(t, game.RoundWinner(), Tie)

	game = Game{
		SideA: &GameSide{
			CloseCombat: ToCards(CardList{
				AllCardsMap["ves"],
			}),
			RangedCombat: &Cards{},
			Siege:        &Cards{},
		},
		SideB: &GameSide{
			CloseCombat:  ToCards(CardList{AllCardsMap["vreemde"]}),
			RangedCombat: &Cards{},
			Siege:        &Cards{},
		},
		WeatherCards: &Cards{},
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
	game := NewGame(&Cards{}, &Cards{}, &Cards{}, &Cards{})
	game.SideA.Passed = true
	game.SideB.Passed = true
	game.NextRound()
	assert.Equal(t, game.Round(), 1)
	assert.Equal(t, game.History[0], Tie)
}

func TestGame_Finished(t *testing.T) {
	game := NewGame(&Cards{}, &Cards{}, &Cards{}, &Cards{})
	assert.False(t, game.Finished())
	game.SideA.Passed = true
	game.SideB.Passed = true
	game.History = []Turn{Tie, Tie}
	assert.True(t, game.Finished())

	game.History = []Turn{PlayerA, PlayerA, Tie}
	t.Log(game.MaxRoundsWon())
	t.Log(game.RoundFinished())
	assert.True(t, game.Finished())

	game.History = []Turn{PlayerA, PlayerB, Tie}
	assert.False(t, game.Finished())

	game.History = []Turn{PlayerA, PlayerB, PlayerA}
	assert.True(t, game.Finished())
}

func TestGame_Winner(t *testing.T) {
	game := Game{
		History: []Turn{PlayerA, PlayerA, Tie},
	}
	t.Log(game.RoundsWon())
	assert.Equal(t, game.Winner(), PlayerA)
}
