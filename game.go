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
func displayBoard(c *gin.Context,
	url string,
	template string,
	choice bool,
	game *gwent.Game,
	player *gwent.GameSide,
	enemy *gwent.GameSide) {
	weather := game.WeatherCards.Effects()
	c.HTML(200, template, gin.H{
		"Weather": gin.H{
			"CardList":     game.WeatherCards,
			"CloseCombat":  weather.Has(gwent.BitingFrost),
			"RangedCombat": weather.Has(gwent.ImpenetrableFog),
			"Siege":        weather.Has(gwent.TorrentialRain),
		},
		"Lives": gin.H{
			"Player": game.LivesLeft(player.Side),
			"Enemy":  game.LivesLeft(enemy.Side),
		},
		"MySide":       player,
		"Player":       player,
		"Enemy":        enemy,
		"DrawHandDeck": player.Hand,
		"Heap":         player.Heap,
		"Url":          url,
		"Choice":       choice,
		"Round":        game.Round(),
	})
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

/*
var game *gwent.Game = gwent.Creategame()

func demoGameHandler(c *gin.Context) {
	game.Sort()
	log.Print(game)
	displayBoard(c, "/move", "game.html", false, game.Player(), game.Enemy())
}



func demoMove(c *gin.Context) {
	row, cardid := getRowAndCard(c)
	card := game.GetCardById(cardid)

	if card != nil {
		log.Println(game.Player().Siege)
		medicEffect := game.PlayMove(card, row, game.Player(), game.Enemy())
		log.Println(game.Player().Siege)
		if medicEffect {
			println("medic !!")
			displayBoard(c, "/move", "table.html", true, game.Player(), game.Enemy())
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
	displayBoard(c, "/move", "game.html", true, game.Player(), game.Enemy())
}
*/
var soloGames = map[string]*gwent.Game{} //cookie->game against IA

func startSoloGame(c *gin.Context) {
	user := c.MustGet("user").(*gwent.PlayerData)
	_sdeckid := c.Param("index")
	deckid, err := strconv.ParseInt(_sdeckid, 10, 64)
	if err != nil {
		panic(err)
	}
	cookie := RandomString(8)
	setcoookie(c, "solo_game", cookie)

	ud := user.Decks.GetByIndex(int(deckid))
	iad := demoData.NewPlayerDeck(gwent.ScoiaTael).Fill()
	sologame, err := gwent.GameFromDecks(ud, iad)
	if err != nil {
		c.String(400, err.Error())
	}
	soloGames[cookie] = sologame
	sologame.Sort()
	c.Redirect(302, "/game")
	//displayBoard(c, "/game", "game.html", false, sologame.Player(), sologame.Enemy())
}

func get_game(c *gin.Context) *gwent.Game {
	cookie, err := c.Cookie("solo_game")
	if err != nil {
		panic(err)
	}
	sologame, ok := soloGames[cookie]
	if !ok {
		//panic("need to start game first")
		return nil
	}
	return sologame
}

func showSoloGame(c *gin.Context) {
	if sologame := get_game(c); sologame != nil {
		displayBoard(c, "/game", "game.html", false, sologame, sologame.Player(), sologame.Enemy())
	} else {
		c.Redirect(302, "/deck/0/")
	}
}

func moveSoloGame(c *gin.Context) {
	sologame := get_game(c)
	player := sologame.SideA
	ia_bot := sologame.SideB
	sologame.Turn = gwent.PlayerA

	row, cardid := getRowAndCard(c)
	if cardid < 0 {
		log.Println("player pass")
		sologame.Pass(player)
	}
	log.Println(sologame.History)
	if sologame.Finished() {
		c.String(200, sologame.Winner().String()+" won !")
		return
	}
	card := sologame.GetCardById(cardid)

	if cardid >= 0 && card == nil {
		c.String(400, "card not found")
	}

	if card != nil {
		medicEffect := sologame.PlayMove(card, row, player, ia_bot)
		if medicEffect {
			displayBoard(c, "/game", "table.html", true, sologame, player, ia_bot)
			return
		}
	}
	sologame.Turn = gwent.PlayerB

	if ia_card, row := gwent.IaMove(sologame, ia_bot, player); ia_card != nil {
		sologame.PlayMove(ia_card, row, ia_bot, player)
	} else {
		sologame.Pass(ia_bot)
	}
	log.Println("lives", sologame.LivesLeft(gwent.PlayerA), sologame.LivesLeft(gwent.PlayerB))

	sologame.Turn = gwent.PlayerA
	displayBoard(c, "/game", "table.html", false, sologame, player, ia_bot)
}
