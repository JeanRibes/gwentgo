package gwent

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGame(t *testing.T) {
	pd := NewPlayerData("jsr022", "mjkjhlkjh")
	deckA := pd.Decks.GetByName("NorthernRealms")
	deckB := pd.Decks.GetByName("Monsters")

	da := []int{0, 57, 56, 3, 8, 39, 9, 1, 20, 28}
	//da := []int{0, 57, 56, 3, 8, 39, 9, 16, 20, 28}
	for _, id := range da {
		deckA.AddToDeck(deckA.Rest.GetById(id))
	}
	db := []int{1, 144, 181, 179, 180, 166, 9, 0, 26, 163}
	//db:=[]int{1,144,181,179,180,166,9,18,26,163}
	for _, id := range db {
		deckB.AddToDeck(deckB.Rest.GetById(id))
	}

	assert.Equal(t, 10, deckA.Deck.Len())
	assert.Equal(t, 10, deckB.Deck.Len())

	game := GameFromDecks(deckA, deckB)
	deckA, deckB = nil, nil

	cardscount := game.Sort()
	assert.Greaterf(t, cardscount, 0, "bad game sorting")
	t.Logf("%d cards in that game", cardscount)

	geralt := game.SideA.Hand.GetByName("geralt-of-rivia")
	game.PlayMove(geralt, geralt.Row, game.SideA, game.SideB)
	if game.SideA.CloseCombat.CountByName("geralt-of-rivia") != 1 {
		t.Fail()
	}
	if card := game.SideA.Hand.GetByName(geralt.Name); card != nil {
		if card.Id == geralt.Id {
			t.Fatal("card not removed from hand")
		}
	}

	ciri := game.SideB.Hand.GetByName("cirilla-fiona-elen-riannon")
	game.PlayMove(ciri, ciri.Row, game.SideB, game.SideA)

	for game.SideA.Hand.Len() > 0 && game.SideB.Hand.Len() > 0 {
		ca := (*game.SideA.Hand)[0]
		game.PlayMove(ca, ca.Row, game.SideA, game.SideB)

		cb := (*game.SideB.Hand)[0]
		game.PlayMove(cb, cb.Row, game.SideB, game.SideA)
	}

	sa, sb := game.Score()
	//t.Log(game)
	assert.Equal(t, 15+15+1+7+8, sa)
	assert.Equal(t, 3*1+1+10+2, sb)

	data, err := json.Marshal(game)
	assert.NoError(t, err, "marshall json")
	assert.NotZero(t, len(data))

	game.Pass(game.SideA)
	assert.False(t, game.RoundFinished())
	assert.Equal(t, 0, game.Round())
	game.Pass(game.SideB)
	assert.Equal(t, 1, game.Round())
	assert.Equal(t, game.History[0], PlayerA)
}
