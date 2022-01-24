package gwent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIaMove(t *testing.T) {
	pd := NewPlayerData("", "")
	deckA := pd.NewPlayerDeck(NorthernRealms)
	deckB := pd.NewPlayerDeck(Monsters)

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
	game := NewGame(&CardList{}, deckA.Deck, &CardList{}, deckB.Deck)
	game.Sort()

	for !game.Finished() {
		assert.NotEqual(t, Tie, game.Turn)
		if card, row := IaMove(game, game.Player(), game.Enemy()); card != nil {
			//t.Logf("(h:%d) move to %s by %s: %s\n",game.Player().Hand.Len(),row,game.Player().Side.String(), card)
			game.PlayMove(card, row, game.Player(), game.Enemy())
		} else {
			game.Pass(game.Player())
		}
		game.Switch()
	}
	t.Log(game.History)
}
