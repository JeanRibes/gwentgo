package main

import (
	"github.com/gin-gonic/gin"
	"gwentgo/gwent"
	"net/http"
)

var demoData = gwent.NewPlayerData("jsr022", "mjkjhlkjh")

func UserDataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user")
		if err != nil {
			c.Set("user", demoData)
			return
			c.String(http.StatusUnauthorized, "unauthorized, no cookie")
			return
		}
		user, ok := userDb[cookie]
		if !ok {
			user = gwent.NewPlayerData("username", cookie)
			userDb[cookie] = user
		}
		c.Set("user", user)
		c.Next()
	}
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	cookie := RandomString(64)
	setcoookie(c, "user", cookie)
	userDb[cookie] = gwent.NewPlayerData(username, cookie)
	c.Redirect(302, "/deck/0/")
}
