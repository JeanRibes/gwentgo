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
	displayBoard(c, "/move", "game.html", false, game.Player(), game.Enemy())
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
	enemy *gwent.GameSide) {
	c.HTML(200, template, gin.H{
		"Weather": game.WeatherCards,
		"MySide":  player,
		"Enemy":   enemy,
		"Hand":    player.Hand,
		"Heap":    player.Heap,
		"Url":     url,
		"Choice":  choice,
		"Round":   game.Round(),
	})
}

func demoMove(c *gin.Context) {
	row, cardid := getRowAndCard(c)
	card := game.GetCardById(cardid)

	if card != nil {
		player := game.Player()
		enemy := game.Enemy()
		medicEffect := game.PlayMove(card, row, player, enemy)
		if medicEffect {
			println("medic !!")
			displayBoard(c, "/move", "table.html", true, player, enemy)
		} else {
			game.Switch()
			displayBoard(c, "/move", "table.html", false, game.Player(), game.Enemy())
		}
	} else {
		c.String(400, "card not found")
	}
}

func demoPass(c *gin.Context) {
	game.Pass(game.Player())
	if game.Finished() {
		log.Printf("rounds won: %s\n", game.History)
		c.String(200, game.Winner().String()+" won !")
		return
	}
	//game.NextRound()
	displayBoard(c, "/move", "table.html", false, game.Player(), game.Enemy())
}

func demoChoice(c *gin.Context) {
	displayBoard(c, "/move", "table.html", true, game.Player(), game.Enemy())
}
