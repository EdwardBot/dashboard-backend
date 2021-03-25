package auth

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//HandleLogin handles the Oauth2 login process
func HandleOAuth(c *gin.Context) {
	authCode, hasCode := c.GetQuery("code")
	if hasCode && authCode != "" {
		client := http.Client{}
		tokenRequest := url.Values{}
		tokenRequest.Set("client_id", os.Getenv("CLIENT_ID"))
		tokenRequest.Set("client_secret", os.Getenv("CLIENT_SECRET"))
		tokenRequest.Set("grant_type", "authorization_code")
		tokenRequest.Set("code", authCode)
		tokenRequest.Set("redirect_uri", os.Getenv("BASE_URL")+"/v1/auth/oauth")
		tokenRequest.Set("scope", "identify email guilds")

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

		res, erro := client.Do(request)
		if erro != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(res.Body)

		var response gin.H
		_ = json.Unmarshal(body, &response)
		if response["error_description"] == "Invalid \"code\" in request." {
			c.JSON(403, gin.H{
				"status": "error",
				"error":  "Code expired.",
			})
			return
		}

		userData, _ := utils.HttpGet("https://discord.com/api/v6/users/@me", map[string]string{
			"Authorization": "Bearer " + response["access_token"].(string),
		})

		sessionCode := rand.Int63()

		user, err := database.FindUser(userData["id"].(string))
		//if err != nil {
		//	go HandleSession(user, response, sessionCode)
		//} else {
		go func(userData gin.H, oauth gin.H, session int64) {
			HandleInitialLogin(userData, oauth)
			HandleSession(user, response, sessionCode)
		}(userData, response, sessionCode)
		//}

		token := jwt.New(jwt.GetSigningMethod("HS256"))

		token.Claims = jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Subject:   userData["id"].(string),
			Id:        strconv.FormatInt(sessionCode, 10),
		}

		tokenStr, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
		if err != nil {
			r := gin.H{
				"status": "error",
				"error":  "Token error",
			}
			b, _ := json.Marshal(r)

			c.HTML(403, "redirect", gin.H{
				"data": string(b),
			})
			log.Println(err)
			return
		}

		r := gin.H{
			"status": "success",
			"token":  tokenStr,
			"id":     userData["id"],
		}

		b, _ := json.Marshal(r)

		c.HTML(200, "redirect", gin.H{
			"data": string(b),
		})
		res.Body.Close()
	} else {
		c.JSON(401, gin.H{
			"status": "Error",
			"error":  "No code sent!",
		})
	}
}

func HandleSession(userId database.User, oauthRes gin.H, code int64) {

}

func HandleInitialLogin(userData gin.H, response gin.H) {
	userId, _ := strconv.ParseInt(userData["id"].(string), 10, 64)
	discriminator, _ := strconv.Atoi(userData["discriminator"].(string))

	guildData, err := utils.HttpGet("https://discord.com/api/v6/users/@me/guilds", map[string]string{
		"Authorization": "Bearer " + response["access_token"].(string),
	})
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Printf("Guilds: %p", guildData)

	userDb := database.User{
		UserID:        userId,
		JoinedAt:      time.Now(),
		UserName:      userData["username"].(string),
		Discriminator: discriminator,
		AvatarID:      userData["avatar"].(string),
	}

	userDb.Save()
}
