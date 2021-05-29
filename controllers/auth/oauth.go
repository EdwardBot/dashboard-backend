package auth

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// HandleOAuth HandleLogin handles the Oauth2 login process
func HandleOAuth(c *gin.Context) {
	authCode, hasCode := c.GetQuery("code")
	if hasCode && authCode != "" {
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

		res, erro := utils.Client.Do(request)
		if erro != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(res.Body)

		var response gin.H
		_ = json.Unmarshal(body, &response)
		if response["error"] != nil {
			c.JSON(403, gin.H{
				"status": "error",
				"error":  "Code expired.",
			})
			return
		}

		userData, _ := utils.HttpGet("https://discord.com/api/v6/users/@me", map[string]string{
			"Authorization": "Bearer " + response["access_token"].(string),
		}, false)

		sessionCode := rand.Int31()

		userID, _ := strconv.ParseInt(userData["id"].(string), 10, 64)

		user, err := database.FindUser(userID)
		go CacheGuilds(userData, response)
		if err != nil {
			go HandleInitialLogin(userData, response)
		}

		if userData["premium_type"] == nil {
			log.Printf("Premium: %s", userData)
		} else {
			user.PremiumType = userData["premium_type"].(int)
			user.Update(user.ID)
		}

		session, err := HandleSession(user, response, sessionCode)

		token := jwt.New(jwt.GetSigningMethod("HS256"))

		token.Claims = jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 603000).Unix(),
			Subject:   userData["id"].(string),
			Id:        strconv.Itoa(int(session.SessionId)),
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
			"status":        "success",
			"token":         tokenStr,
			"id":            userData["id"],
			"username":      userData["username"],
			"discriminator": userData["discriminator"],
			"session_id":    session.SessionId,
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

func HandleSession(user database.User, oauthRes gin.H, code int32) (database.Session, error) {
	_, err := database.FindSession(code)
	if err == nil {
		log.Println(err)
	}
	session := database.Session{
		ID:           primitive.NewObjectID(),
		UserID:       user.UserID,
		AccessToken:  oauthRes["access_token"].(string),
		RefreshToken: oauthRes["refresh_token"].(string),
		RefreshedAt:  time.Now(),
		ExpiresIn:    604000,
		SessionId:    code,
	}
	err = session.Save()
	return session, err
}

func HandleInitialLogin(userData gin.H, response gin.H) {
	userId, _ := strconv.ParseInt(userData["id"].(string), 10, 64)
	discriminator, _ := userData["discriminator"].(string)

	guildData, err := utils.HttpGet("https://discord.com/api/v6/users/@me/guilds", map[string]string{
		"Authorization": "Bearer " + response["access_token"].(string),
	}, true)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	guilds := guildData["guilds"].([]interface{})

	guildIDs := make([]int64, 0, 100)

	log.Printf("Caching guilds for user \"%s\"", userData["username"].(string))
	for g := range guilds {
		guild := guilds[g].(map[string]interface{})
		dId, _ := strconv.ParseInt(guild["id"].(string), 10, 64)
		var icon string
		if guild["icon"] == nil {
			icon = "default"
		} else {
			icon = guild["icon"].(string)
		}
		_, hasBot := database.FindGConf(guild["id"].(string))
		if hasBot == nil {
			toSave := database.Guild{
				GuildID:    dId,
				Name:       guild["name"].(string),
				Icon:       icon,
				HasBot:     hasBot == nil,
				OwnerId:    -1,
				HasPremium: false,
				ID:         primitive.NewObjectID(),
			}
			tryCacheGuild(toSave)
			guildIDs = append(guildIDs, dId)
		}
	}
	log.Printf("Caching guilds done for user \"%s\"", userData["username"].(string))

	userDb := database.User{
		ID:            primitive.NewObjectID(),
		UserID:        userId,
		JoinedAt:      time.Now(),
		UserName:      userData["username"].(string),
		Discriminator: discriminator,
		AvatarID:      userData["avatar"].(string),
		Guilds:        guildIDs,
	}

	userDb.Save()
}

func CacheGuilds(userData, response gin.H) {

}

func tryCacheGuild(guild database.Guild) {
	g, err := database.FindGuild(guild.GuildID)
	if err != nil {
		guild.Save()
	} else {
		guild.Update(g.ID)
	}
}
