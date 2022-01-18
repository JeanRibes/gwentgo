package main

import (
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"log"
	"strconv"
)

func gameHandler(c *gin.Context) {
	cookie, err0 := c.Cookie("gamecookie")
	side, err1 := c.Cookie("side")
	if err0 != nil || err1 != nil {
		c.String(400, "error"+err0.Error()+err1.Error())
		return
	}
	data := games[cookie]
	var game *gwent.Game = data.Game
	var player *gwent.GameSide
	var enemy *gwent.GameSide
	if side == "A" {
		player = game.SideA
		enemy = game.SideB
	} else {
		enemy = game.SideA
		player = game.SideB
	}
	c.HTML(200, "game.html", gin.H{
		"Weather": game.WeatherCards,
		"MySide":  player,
		"Enemy":   enemy,
	})
}

var game *gwent.Game = gwent.Creategame()

func demoGameHandler(c *gin.Context) {
	game.Sort()
	log.Print(game)
	displayBoard(c, "/move", "game.html", false, game.SideA, game.SideB, game.SideA.Hand, game.SideA.Heap)
}

func getRowAndCard(c *gin.Context) (row gwent.Row, cardid int) {
	sid := c.PostForm("id")
	srow := c.PostForm("row")
	if sid == "" || srow == "" {
		c.String(400, "bad form, missing number 'id'")
		return
	}
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(400, "not a number")
		return
	}
	return gwent.RowFromString(srow), int(id)
}

func displayBoard(c *gin.Context,
	url string,
	template string,
	choice bool,
	player *gwent.GameSide,
	enemy *gwent.GameSide,
	hand *gwent.Cards,
	heap *gwent.Cards) {
	c.HTML(200, template, gin.H{
		"Weather": game.WeatherCards,
		"MySide":  player,
		"Enemy":   enemy,
		"Hand":    hand,
		"Heap":    heap,
		"Url":     url,
		"Choice":  choice,
	})
}

func demoMove(c *gin.Context) {
	row, cardid := getRowAndCard(c)
	card := game.GetCardById(cardid)
	if card != nil {
		medicEffect := game.PlayMove(card, row, game.SideA, game.SideB)
		if medicEffect {
			println("medic !!")
			displayBoard(c, "/demo", "table.html", true, game.SideA, game.SideB, game.SideA.Hand, game.SideA.Heap)
		} else {
			displayBoard(c, "/demo", "table.html", false, game.SideA, game.SideB, game.SideA.Hand, game.SideA.Heap)
		}
	} else {
		c.String(400, "card not found")
	}
}

/*
func demoChoice(c *gin.Context) {
	row, cardid := getRowAndCard(c)
	card := game.GetCardById(cardid)

}*/
