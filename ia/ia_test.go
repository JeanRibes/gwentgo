package ia

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gwentgo/gwent"
	"testing"
)

func TestIaMove(t *testing.T) {
	pd := gwent.NewPlayerData("", "")
	deckA := pd.NewPlayerDeck(gwent.NorthernRealms)
	deckB := pd.NewPlayerDeck(gwent.Monsters)

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
	game := gwent.NewGame(&gwent.CardList{}, deckA.Deck, &gwent.CardList{}, deckB.Deck)
	game.Sort()

	for !game.Finished() {
		assert.NotEqual(t, gwent.Tie, game.Turn)
		if card, row := IaMove(game, game.Player(), game.Enemy()); card != nil {
			//t.Logf("(h:%d) move to %s by %s: %s\n",game.Player().Hand.Len(),row,game.Player().Side.String(), card)
			game.PlayMove(card, row, game.Player(), game.Enemy())
		} else {
			game.Pass(game.Player())
		}
		game.Switch()
	}
	t.Log(game.History)
	/*for _, rec := range game.Recording {
		t.Log(rec)
	}*/

	data, err := json.Marshal(game.Recording)
	assert.NoError(t, err)
	t.Log(string(data))
}
