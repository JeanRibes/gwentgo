package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"strconv"
)

func showDeck(c *gin.Context) {
	userData := c.MustGet("user").(*gwent.PlayerData)
	_index, _ := strconv.ParseInt(c.Param("index"), 10, 64)
	index := int(_index)

	deck := userData.Decks.GetByIndex(index)
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
		"PrevDeck": ((index - 1 + len(*userData.Decks)) % len(*userData.Decks)),
		"NextDeck": (index + 1) % len(*userData.Decks),
	})
}

func addToDeck(c *gin.Context) {
	userData := c.MustGet("user").(*gwent.PlayerData)
	sid := c.Param("id")
	sindex := c.Param("index")
	_id, _ := strconv.ParseInt(sid, 10, 64)
	_index, _ := strconv.ParseInt(sindex, 10, 64)
	id := int(_id)
	index := int(_index)
	deck := userData.Decks.GetByIndex(index)
	card := deck.Rest.GetById(id)
	if card == nil {
		c.String(400, "card not found")
		return
	}
	deck.AddToDeck(card)
	c.String(200, "ok, added")
}

func removeFromDeck(c *gin.Context) {
	userData := c.MustGet("user").(*gwent.PlayerData)
	sid := c.Param("id")
	sindex := c.Param("index")
	_id, _ := strconv.ParseInt(sid, 10, 64)
	_index, _ := strconv.ParseInt(sindex, 10, 64)
	id := int(_id)
	index := int(_index)
	deck := userData.Decks.GetByIndex(index)
	card := deck.Deck.GetById(id)
	if card == nil {
		c.String(400, "card not found")
		return
	}
	deck.RemoveFromDeck(card)
	c.String(200, "ok, removed")
}

func editDeck(c *gin.Context) {
	userData := c.MustGet("user").(*gwent.PlayerData)
	_index, _ := strconv.ParseInt(c.Param("index"), 10, 64)
	index := int(_index)

	deck := userData.Decks.GetByIndex(index)

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
		(*userData.Decks)[index] = userData.NewPlayerDeck(newFaction)
	}
	showDeck(c)
}
func newDeck(c *gin.Context) {
	userData := c.MustGet("user").(*gwent.PlayerData)
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
	*userData.Decks = append(*userData.Decks, userData.NewPlayerDeck(faction))
	c.Redirect(302, fmt.Sprintf("/deck/%d", len(*userData.Decks)-1))
}
