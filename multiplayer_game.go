package main

import (
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"io"
	"net/http"
	"sync"
)

type MultiplayerGame struct {
	Key            string
	WaitingChannel chan struct{}
	sync.Mutex

	Game  *gwent.Game
	ChanA chan bool
	UserA *gwent.PlayerData
	UserB *gwent.PlayerData
	ChanB chan bool
}

var multiplayer_games = map[string]*MultiplayerGame{}

// createMPgame url: /multi/create
func createMPgame(c *gin.Context) {
	key := RandomString(32)
	multiplayer_games[key] = &MultiplayerGame{
		Key:            key,
		WaitingChannel: make(chan struct{}, 1),
	}
	setcoookie(c, "gamekey", key)
	c.Redirect(http.StatusFound, "/multi/join")
}

// joinGame url: /multi/join/:key (key optionnel)
func joinGame(c *gin.Context) {
	user := c.MustGet("user").(*gwent.PlayerData)
	var key string
	if key = c.Param("gamekey"); key != "" {
		setcoookie(c, "gamekey", key)
	} else {
		key, _ = c.Cookie("gamekey")
	}
	if key == "" {
		c.String(400, "no game key")
		return
	}
	lobby, ok := multiplayer_games[key]
	if !ok {
		c.String(400, "no game found")
		return
	}
	if c.Param("gamekey") != "" { //joined via link
		lobby.WaitingChannel <- struct{}{}
	}
	side := lobby.AssignUser(user)
	setcoookie(c, "side", side.String())
	c.HTML(200, "waiting.html", gin.H{"Key": key})
}

// waitingRoom url: /multi/wait
func waitingRoom(c *gin.Context) {
	cookie, err := c.Cookie("gamekey")
	if err != nil || cookie == "" {
		c.String(400, "missing game cookie")
		return
	}
	lobby, ok := multiplayer_games[cookie]
	if !ok {
		c.String(400, "no game found")
		return
	}

	c.Stream(func(w io.Writer) bool {
		<-lobby.WaitingChannel
		c.SSEvent("ready", "opponent joined !")
		lobby.WaitingChannel <- struct{}{} //put back into chan for the other side
		return false
	})
	return
}

func (mg *MultiplayerGame) IsFull() bool {
	mg.Lock()
	val := mg.UserA != nil && mg.UserB != nil
	mg.Unlock()
	return val
}

func (mg *MultiplayerGame) AssignUser(user *gwent.PlayerData) gwent.Turn {
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
