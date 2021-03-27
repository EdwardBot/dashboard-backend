package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func HasAuth(c *gin.Context) {
	if c.Request.Header.Get("Authorization") == "" {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "No token provided",
		})
	}
	token := strings.Replace(c.Request.Header.Get("Authorization"), "Bearer ", "", 1)
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  "Internal error",
		})
	}
	if !jwtToken.Valid {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Invalid token",
		})
	}
	tokenData := jwtToken.Claims.(jwt.MapClaims)
	log.Println(tokenData)
	exp := tokenData["exp"].(float64)
	now, _ := strconv.ParseFloat(strconv.FormatInt(time.Now().Unix(), 10), 64)
	if exp < now {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Expired token",
		})
	}
	sessionId, _ := strconv.ParseInt(tokenData["jti"].(string), 10, 32)
	session, err := database.FindSession(int32(sessionId))
	if err != nil {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Invalid session",
		})
	}
	userId, _ := strconv.ParseInt(tokenData["sub"].(string), 10, 64)
	if userId != session.UserID {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Invalid token",
		})
	}
}
