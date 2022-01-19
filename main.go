package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"log"
	"os"
	"os/signal"
)

type Lobby struct {
	Cookie string
	Name   string
	Key    string
}
type Game struct {
	Name   string
	Cookie string
	Game   *gwent.Game
}

type Db struct {
	Lobbies map[string]Lobby
	Games   map[string]Game
}

var lobbies map[string]Lobby //key -> cookie

var games map[string]Game //cookie -> game

func init() {
	lobbies = map[string]Lobby{}
	games = map[string]Game{}
}

func setcoookie(c *gin.Context, key string, value string) {
	c.SetCookie(key, value, 3600, "/", "localhost", false, false)
}

func waitingRoom(c *gin.Context, key string) {
	lobby, exists := lobbies[key]
	if exists {
		if _, gexists := games[lobby.Cookie]; gexists {
			setcoookie(c, "side", "A")
			delete(lobbies, key)
			c.Redirect(302, "/game")
		} else {
			c.HTML(200, "waiting.html", &lobby)
		}
	} else {
		c.String(200, "lobby not found, please create another")
	}
}

func main() {
	load()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		save()
		os.Exit(0)
	}()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./assets")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", &struct {
			Lobbies map[string]Lobby
		}{Lobbies: lobbies})
	})
	r.GET("/lobby", func(c *gin.Context) {
		s := `<section id="Lobbies" hx-swap-oob="true">
<button hx-get="/lobby">Refresh Lobbies</button>
<ul>`
		for _, lobby := range lobbies {
			s += `<li><a href="/join/` + lobby.Key + `">` + lobby.Name + "</a></li>\n"
		}
		c.String(200, s+"</ul></section>")
	})

	r.GET("/join/:key", func(c *gin.Context) {
		key := c.Param("key")
		lobby, exists := lobbies[key]
		if exists {
			setcoookie(c, "gamecookie", lobby.Cookie)
			games[lobby.Cookie] = Game{
				Name:   lobby.Name,
				Cookie: lobby.Cookie,
				Game:   gwent.Creategame(),
			}
			setcoookie(c, "side", "B")
			c.Redirect(302, "/game")
		} else {
			c.Redirect(302, "/")
		}
	})

	r.POST("/create", func(c *gin.Context) {
		key := RandomString(16)
		cookie := RandomString(16)
		lobbies[key] = Lobby{
			Key:    key,
			Cookie: cookie,
			Name:   c.PostForm("name"),
		} //  SetCookie(name string, value string, maxAge int, path string, domain string, secure
		setcoookie(c, "gamecookie", cookie)
		setcoookie(c, "gamekey", key)
		//c.JSON(201, Lobbies[key])
		waitingRoom(c, key)
	})

	r.GET("/info", func(c *gin.Context) {
		cookie, err := c.Cookie("gamecookie")
		if cookie == "" || err != nil {
			c.String(400, "plz join a game")
			return
		}
		game, exists := games[cookie]
		if exists {
			c.String(200, fmt.Sprintf("game %s", game.Name))
		} else {
			c.String(200, "plz join a game")
		}
	})

	r.GET("/waiting", func(c *gin.Context) {
		cookie, err := c.Cookie("gamekey")
		if cookie == "" || err != nil {
			log.Print(cookie, err)
			c.String(400, "plz join a game")
			return
		}
		waitingRoom(c, cookie)
	})
	r.GET("/save", func(c *gin.Context) {
		save()
		c.String(200, "saved !")
	})

	//r.GET("/game", gameHandler)
	r.GET("/demo", demoGameHandler)
	r.POST("/move", demoMove)
	r.POST("/pass", demoPass)

	r.GET("/choicedemo", demoChoice)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func save() {
	f, err := os.OpenFile("db.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(&Db{
		Lobbies: lobbies,
		Games:   games,
	})
	_, err = f.Write(data)
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func load() {
	f, err := os.Open("db.json")
	if err != nil {
		return
	}
	var db Db
	if err := json.NewDecoder(f).Decode(&db); err != nil {
		panic(err)
	}
	if db.Games != nil {
		games = db.Games
	}
	if db.Lobbies != nil {
		lobbies = db.Lobbies
	}
}
