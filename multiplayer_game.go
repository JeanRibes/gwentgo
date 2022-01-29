package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"io"
	"net/http"
	"strconv"
	"sync"
)

type MultiplayerGame struct {
	Key            string
	WaitingChannel chan struct{}
	sync.Mutex
	Game *gwent.Game

	ChanA chan bool
	UserA *gwent.PlayerData
	DeckA *gwent.PlayerDeck

	ChanB chan bool
	UserB *gwent.PlayerData
	DeckB *gwent.PlayerDeck
}

var multiplayer_games = map[string]*MultiplayerGame{}

// createMPgame url: /multi/create
func createMPgame(c *gin.Context) {
	key := RandomString(32)
	multiplayer_games[key] = &MultiplayerGame{
		Key:            key,
		WaitingChannel: make(chan struct{}, 1),
	}
	setcoookie(c, GKEY, key)
	delcookie(c, SIDE)
	c.Redirect(http.StatusFound, "/multi/join")
}

const MP = "multiplayer_game"
const GKEY = "gamekey"

func MultiGameMiddleware(c *gin.Context) {
	cookie, err := c.Cookie(GKEY)
	if err != nil || cookie == "" {
		c.String(400, "missing game cookie")
		c.Abort()
		return
	}
	lobby, ok := multiplayer_games[cookie]
	if !ok {
		c.String(400, "no game found")
		c.Abort()
		return
	}

	if _side, err := c.Cookie(SIDE); err == nil {
		side := gwent.UnstringTurn(_side)
		if side != gwent.Tie {
			c.Set(SIDE, side)
		}
	}

	c.Set(MP, lobby)
	c.Next()
}

const SIDE = "side"

// joinGameP1 for the first player, who initiated the lobby
func joinGameP1(c *gin.Context) {
	user := c.MustGet("user").(*gwent.PlayerData)
	key, _ := c.Cookie(GKEY)
	lobby, ok := multiplayer_games[key]
	if !ok {
		c.String(400, "no game found, bad key")
		return
	}
	if side, _ := c.Cookie(SIDE); side == "" {
		setcoookie(c, SIDE, lobby.AssignUser(user).String())
	}
	c.HTML(200, "waiting.html", gin.H{"Key": key, "User": user})
}

func joinGameP2(c *gin.Context) {
	user := c.MustGet("user").(*gwent.PlayerData)
	key := c.Param(GKEY)
	lobby, ok := multiplayer_games[key]
	if !ok {
		c.String(400, "no game found")
		return
	}

	if ckey, _ := c.Cookie(GKEY); ckey != key { // first time visiting this URL
		setcoookie(c, GKEY, key)
		lobby.WaitingChannel <- struct{}{}
		setcoookie(c, SIDE, lobby.AssignUser(user).String())
	}
	c.HTML(200, "waiting.html", gin.H{"Key": key, "User": user})
}

// waitingRoom url: /multi/wait
func waitingRoom(c *gin.Context) {
	lobby := c.MustGet(MP).(*MultiplayerGame)

	c.Stream(func(w io.Writer) bool {
		<-lobby.WaitingChannel
		c.SSEvent("ready", "opponent joined !")
		lobby.WaitingChannel <- struct{}{} //put back into chan for the other side
		return false
	})
	return
}

func chooseDeck(c *gin.Context) {
	user := c.MustGet(USER).(*gwent.PlayerData)
	lobby := c.MustGet(MP).(*MultiplayerGame)
	side := c.MustGet(SIDE).(gwent.Turn)

	sindex := c.Param("index")
	index, err := strconv.ParseInt(sindex, 10, 64)
	if err != nil {
		panic(err)
		return
	}
	deck := user.Decks.GetByIndex(int(index))
	if deck == nil {
		panic("invalid index")
		return
	}
	lobby.SetDeck(deck, side)
	c.String(200, "you chose %s", deck.Name)
}

// startGame is called by the two parties
func startGame(c *gin.Context) {
	lobby := c.MustGet(MP).(*MultiplayerGame)
	side := c.MustGet(SIDE).(gwent.Turn)

	if !lobby.Started() {
		if (side == gwent.PlayerB && lobby.DeckA == nil) || (side == gwent.PlayerA && lobby.DeckB == nil) {
			c.String(200, "Enemy is not ready, they must choose their deck")
		}
		if err := lobby.StartGame(); err != nil {
			c.String(400, err.Error())
			return
		}
	}
	c.String(201, "ok")
}

func (mg *MultiplayerGame) IsFull() bool {
	mg.Lock()
	val := mg.UserA != nil && mg.UserB != nil
	mg.Unlock()
	return val
}

func (mg *MultiplayerGame) AssignUser(user *gwent.PlayerData) gwent.Turn {
	mg.Lock()
	defer mg.Unlock()
	if mg.UserA == nil {
		mg.UserA = user
		return gwent.PlayerA
	}
	if mg.UserB == nil {
		mg.UserB = user
		return gwent.PlayerB
	}
	return gwent.Tie
}

func (mg *MultiplayerGame) SetDeck(deck *gwent.PlayerDeck, side gwent.Turn) {
	mg.Lock()
	if side == gwent.PlayerA {
		mg.DeckA = deck
	}
	if side == gwent.PlayerB {
		mg.DeckB = deck
	}
	mg.Unlock()
}

func (mg *MultiplayerGame) StartGame() error {
	mg.Lock()
	defer mg.Unlock()
	if mg.DeckA != nil && mg.DeckB != nil && mg.Game == nil {
		game, err := gwent.GameFromDecks(mg.DeckA, mg.DeckB)
		if err != nil {
			return err
		}
		mg.Game = game
		return nil
	}
	return fmt.Errorf("invalid condition")
}

func (mg *MultiplayerGame) Started() bool {
	mg.Lock()
	defer mg.Unlock()
	return mg.Game != nil
}

// GameReady returns if both players are connected and have chosen their decks, the game is ready to start
func (mg *MultiplayerGame) GameReady() bool {
	mg.Lock()
	defer mg.Unlock()
	return mg.IsFull() && mg.DeckA != nil && mg.DeckB != nil
}
