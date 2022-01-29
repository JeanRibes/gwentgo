package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
)

func setcoookie(c *gin.Context, key string, value string) {
	c.SetCookie(key, value, 360000000, "/", "localhost", false, false)
}
func delcookie(c *gin.Context, key string) {
	c.SetCookie(key, "", -1, "/", "localhost", false, false)
}

func main() {
	println("Loading data...")
	loadData()
	loadmpg()
	println("done")
	go backupRoutine()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		saveData()
		savempg()
		os.Exit(0)
	}()
	logger := gin.Logger()
	r := gin.New()
	//r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/save", logger, func(c *gin.Context) {
		saveData()
		c.String(200, "saved !")
	})

	/*r.GET("/sse", ssePage)
	r.GET("/sse/event", sseTest)
	r.GET("/sse/counter", dummyreplace)*/

	r.Group("/multi").
		Use(logger, UserDataMiddleware()).
		GET("/create", createMPgame).
		GET("/join", joinGameP1).
		GET("/join/:gamekey", joinGameP2).
		Use(MultiGameMiddleware).
		GET("/wait", waitingRoom).
		GET("/choosedeck/:index", chooseDeck).
		GET("/game/event", gameEventListener).
		GET("/game", gameHandler).
		POST("/game", moveHandler)

	//r.GET("/game", gameHandler)

	/*r.GET("/demo", demoGameHandler)
	r.POST("/move", demoMove)
	r.POST("/pass", demoPass)

	r.GET("/choicedemo", demoChoice)*/

	r.GET("/register", logger, func(c *gin.Context) {
		c.File("templates/register.html")
	})
	r.POST("/register", logger, register)

	r.Group("/deck").
		Use(logger, UserDataMiddleware()).
		GET("/:index/", showDeck).
		POST("/:index/", editDeck).
		POST("/", newDeck).
		GET("/:index/start", startSoloGame).
		POST("/:index/add/:id", addToDeck).
		POST("/:index/remove/:id", removeFromDeck)

	r.POST("/game", logger, moveSoloGame)
	r.GET("/game", logger, showSoloGame)

	println("starting server")
	//r.Run() // listen and serve on 0.0.0.0:8080
	/*go log.Println(http3.ListenAndServe(":443",
		"/home/jean/dev/gemini-server/reverseproxy/certs/crt.pem",
		"/home/jean/dev/gemini-server/reverseproxy/certs/key.pem", r))
	go func() {
		log.Println(http.ListenAndServe(":80", r))
		log.Println("http stopped")
	}()*/
	r.Run(":8080")
}
