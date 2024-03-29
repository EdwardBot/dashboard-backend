package auth

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func HandleRefresh(c *gin.Context) {
	if !c.MustGet("hasAuth").(bool) {
		return
	}
	data := c.MustGet("body").(map[string]interface{})

	if data == nil {
		c.JSON(401, gin.H{
			"status": "error",
		})
		return
	}

	sessionId, _ := strconv.ParseInt(data["id"].(string), 10, 32)

	var session database.Session
	r := database.Conn.First(&session, sessionId)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(401, gin.H{
			"status": "error",
			"error":  "Invalid session",
		})
		return
	}

	tokenRequest := url.Values{}
	tokenRequest.Set("client_id", os.Getenv("CLIENT_ID"))
	tokenRequest.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	tokenRequest.Set("grant_type", "refresh_token")
	tokenRequest.Set("refresh_token", session.RefreshToken)
	tokenRequest.Set("redirect_uri", os.Getenv("BASE_URL")+"/v1/auth/oauth")
	tokenRequest.Set("scope", "identify guilds")

	request, err := http.NewRequest("POST", "https://discord.com/api/v8/oauth2/token", strings.NewReader(tokenRequest.Encode()))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(tokenRequest.Encode())))

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "Error",
			"error":  "Error getting token!",
		})
		return
	}

	res, erro := utils.Client.Do(request)
	if erro != nil {
		log.Fatal(err)
	}

	resBody, _ := ioutil.ReadAll(res.Body)

	var response gin.H
	_ = json.Unmarshal(resBody, &response)

	if response["error"] == nil {
		session.RefreshToken = response["refresh_token"].(string)
		session.AccessToken = response["access_token"].(string)
		session.RefreshedAt = time.Now().Unix()
		database.Conn.Save(&session)
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))

	userId := strconv.FormatInt(session.UserID, 10)

	token.Claims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 603000).Unix(),
		Subject:   userId,
		Id:        strconv.Itoa(int(session.ID)),
	}

	tokenStr, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		c.JSON(403, gin.H{
			"status": "error",
			"error":  "Token error",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":    "success",
		"token":     tokenStr,
		"sessionId": session.ID,
		"id":        userId,
	})
}
