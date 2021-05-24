package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

func ParseToken(c *gin.Context) jwt.MapClaims {
	token := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1)
	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if jwtToken == nil {
		log.Println("FIXME! Token is nil, but it shouldn't be nil!")
		return nil
	}
	tokenData := jwtToken.Claims
	return tokenData.(jwt.MapClaims)
}
