package main

import (
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
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

func demoGameHandler(c *gin.Context) {
	game := gwent.Creategame()
	c.HTML(200, "game.html", gin.H{
		"Weather": game.WeatherCards,
		"MySide":  game.SideB,
		"Enemy":   game.SideA,
	})
}
