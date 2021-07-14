package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
	"time"
)

func HasAuth(c *gin.Context) {
	c.Set("hasAuth", false)
	if c.Request.Header.Get("Authorization") == "" {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "No token provided",
		})
		return
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
		return
	}
	if !jwtToken.Valid {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Invalid token",
		})
		return
	}
	tokenData := jwtToken.Claims.(jwt.MapClaims)

	exp := tokenData["exp"].(float64)
	now, _ := strconv.ParseFloat(strconv.FormatInt(time.Now().Unix(), 10), 64)
	if exp < now {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Expired token",
		})
		return
	}
	sessionId, _ := strconv.ParseInt(tokenData["jti"].(string), 10, 32)
	var session database.Session
	r := database.Conn.First(&session, sessionId)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Invalid session",
		})
		return
	}
	userId, _ := strconv.ParseInt(tokenData["sub"].(string), 10, 64)
	if userId != session.UserID {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Invalid token",
		})
		return
	}
	c.Set("hasAuth", true)
	c.Set("userId", tokenData["sub"].(string))
	c.Next()
}
