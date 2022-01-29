package main

import (
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"net/http"
)

var demoData = gwent.NewPlayerData("demouser", "mjkjhlkjh")

func init() {
	for _, deck := range *demoData.Decks {
		deck.Name = "Demo deck !!!"
	}
}

const USER = "user"

func UserDataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(USER)
		if err != nil {
			c.String(http.StatusUnauthorized, "unauthorized, no cookie")
			c.Abort()
			return
		}
		user, ok := userDb[cookie]
		if !ok {
			user = gwent.NewPlayerData("username", cookie)
			userDb[cookie] = user
		}
		c.Set(USER, user)

		c.Next()
	}
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	cookie := RandomString(64)
	setcoookie(c, USER, cookie)
	userDb[cookie] = gwent.NewPlayerData(username, cookie)
	c.Redirect(302, "/deck/0/")
}
