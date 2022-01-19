package gwent

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func aTestGame(t *testing.T) {
	game := Creategame()
	cardscount := game.Sort()
	assert.Greaterf(t, cardscount, 0, "bad game sorting")
	t.Logf("%d cards in that game", cardscount)

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
		t.Log(game)
		t.Fatal("Muster didn't work")
	}
	//t.Log(game)

	data, err := json.Marshal(game)
	assert.NoError(t, err, "marshel json")
	assert.NotZero(t, len(data))
	//t.Log(string(data))

	assert.Truef(t, game.Check(), "game init check")
	assert.Equal(t, game.ScoreA(), 15)
	assert.Equal(t, game.SideA.CloseCombat.GetByName("geralt-of-rivia").score, 15)
	assert.Equal(t, game.ScoreB(), 9)

	troiseiem := game.SideA.Hand.GetByName("villentretenmerth")
	game.PlayMove(troiseiem, troiseiem.Row, game.SideA, game.SideB)

	merged := game.Merge()
	assert.Equal(t, len(merged), 5)

	quatrieeme := game.SideB.Hand.GetByName("scorch")
	game.PlayMove(quatrieeme, quatrieeme.Row, game.SideB, game.SideA)
	assert.Equal(t, game.ScoreA(), 4)

	cinquieme := game.SideA.Hand.GetByName("blue-stripes-commando")
	game.PlayMove(cinquieme, cinquieme.Row, game.SideA, game.SideB)
	assert.Equal(t, game.ScoreA(), 16)

	sixieme := game.SideB.Hand.GetByName("biting-frost")
	game.PlayMove(sixieme, Weather, game.SideB, game.SideA)

	assert.Truef(t, game.WeatherCards.Effects().Has(BitingFrost), "no biting frost")

	assert.Equal(t, game.ScoreA(), 4)
	assert.Equal(t, game.ScoreB(), 3)

	septieme := game.SideA.Hand.GetByName(AllCards[9].Name)
	assert.Equal(t, game.SideA.Hand.Len(), 7)
	game.PlayMove(septieme, septieme.Row, game.SideA, game.SideB)
	assert.Equal(t, game.SideA.Hand.Len(), 8)
	assert.Equal(t, game.ScoreB(), 3)

	huitieme := game.SideB.Hand.GetByName("saesenthessis")
	game.PlayMove(huitieme, RangedCombat, game.SideB, game.SideA)
	assert.Equal(t, game.ScoreB(), 13)
}
