package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
)

func HandleUserInfo(c *gin.Context) {
	token := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1)
	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	tokenData := jwtToken.Claims
	uIdS := tokenData.(jwt.MapClaims)["sub"]
	uId, _ := strconv.ParseInt(uIdS.(string), 10, 64)
	user, e := database.FindUser(uId)
	if e != nil {
		c.JSON(404, gin.H {
			"status": "error",
			"error": "no user found",
		})
		return
	}
	c.JSON(200, gin.H {
		"username": user.UserName,
		"discriminator": user.Discriminator,
	})
}
