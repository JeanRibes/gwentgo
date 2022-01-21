package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"strconv"
)

var demoData = gwent.NewPlayerData("jsr022", "mjkjhlkjh")

func showDeck(c *gin.Context) {
	_index, _ := strconv.ParseInt(c.Param("index"), 10, 64)
	index := int(_index)

	deck := demoData.Decks.GetByIndex(index)
	if !deck.Deck.CheckNil() || !deck.Rest.CheckNil() {
		c.String(400, "server error: deck is nil")
		return
	}
	c.HTML(200, "deck.html", gin.H{
		"Deck":     deck.Deck,
		"Name":     deck.Name,
		"Faction":  deck.Faction().String(),
		"Leader":   deck.Leader,
		"Rest":     deck.Rest, //.FilterFaction(deck.Faction()).Removes(*deck.Deck),
		"PrevDeck": ((index - 1 + len(*demoData.Decks)) % len(*demoData.Decks)),
		"NextDeck": (index + 1) % len(*demoData.Decks),
	})
}

func addToDeck(c *gin.Context) {
	sid := c.Param("id")
	sindex := c.Param("index")
	_id, _ := strconv.ParseInt(sid, 10, 64)
	_index, _ := strconv.ParseInt(sindex, 10, 64)
	id := int(_id)
	index := int(_index)
	deck := demoData.Decks.GetByIndex(index)
	card := deck.Rest.GetById(id)
	if card == nil {
		c.String(400, "card not found")
		return
	}
	deck.AddToDeck(card)
	c.String(200, "ok, added")
}

func removeFromDeck(c *gin.Context) {
	sid := c.Param("id")
	sindex := c.Param("index")
	_id, _ := strconv.ParseInt(sid, 10, 64)
	_index, _ := strconv.ParseInt(sindex, 10, 64)
	id := int(_id)
	index := int(_index)
	deck := demoData.Decks.GetByIndex(index)
	card := deck.Deck.GetById(id)
	if card == nil {
		c.String(400, "card not found")
		return
	}
	deck.RemoveFromDeck(card)
	c.String(200, "ok, removed")
}

func editDeck(c *gin.Context) {
	_index, _ := strconv.ParseInt(c.Param("index"), 10, 64)
	index := int(_index)

	deck := demoData.Decks.GetByIndex(index)

	name := c.PostForm("name")
	if name != "" {
		deck.Name = name
	}
	if snewFaction := c.PostForm("faction"); snewFaction != "" {
		newFaction, err := gwent.FactionFromString(snewFaction)
		if err != nil {
			c.String(400, "invalid faction")
			return
		}
		(*demoData.Decks)[index] = demoData.NewPlayerDeck(newFaction)
	}
	showDeck(c)
}
func newDeck(c *gin.Context) {
	sfaction := c.PostForm("faction")
	if sfaction == "" {
		c.String(400, "need key 'faction' (string)")
		return
	}
	faction, err := gwent.FactionFromString(sfaction)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	*demoData.Decks = append(*demoData.Decks, demoData.NewPlayerDeck(faction))
	c.Redirect(302, fmt.Sprintf("/deck/%d", len(*demoData.Decks)-1))
}
