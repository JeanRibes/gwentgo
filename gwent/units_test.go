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
