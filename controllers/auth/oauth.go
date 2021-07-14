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

// HandleOAuth handles the Oauth2 login process
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
			c.JSON(http.StatusInternalServerError, gin.H{
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

		userID, _ := strconv.ParseInt(userData["id"].(string), 10, 64)

		var user database.User
		r := database.Conn.Model(&database.User{}).First(&user, userID)
		go HandleInitialLogin(userData, response, errors.Is(r.Error, gorm.ErrRecordNotFound), user)

		if userData["premium_type"] == nil {
			log.Printf("Premium: %s", userData)
		} else {
			user.PremiumType = int(userData["premium_type"].(float64))
			database.Conn.Model(&database.User{}).Where("user_id = ?", user.UserID).Save(&user)
		}

		session, err := HandleSession(user, response)

		token := jwt.New(jwt.GetSigningMethod("HS256"))

		token.Claims = jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 603000).Unix(),
			Subject:   userData["id"].(string),
			Id:        strconv.FormatInt(int64(session.ID), 10),
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

		rData := gin.H{
			"status":        "success",
			"token":         tokenStr,
			"id":            userData["id"],
			"username":      userData["username"],
			"discriminator": userData["discriminator"],
			"session_id":    session.ID,
		}

		b, _ := json.Marshal(rData)

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

func HandleSession(user database.User, oauthRes gin.H) (database.Session, error) {
	session := database.Session{
		UserID:       user.UserID,
		AccessToken:  oauthRes["access_token"].(string),
		RefreshToken: oauthRes["refresh_token"].(string),
		RefreshedAt:  time.Now().Unix(),
		ExpiresIn:    604000,
	}
	r := database.Conn.Create(&session)
	return session, r.Error
}

func HandleInitialLogin(userData gin.H, response gin.H, initial bool, user database.User) {
	userId, _ := strconv.ParseInt(userData["id"].(string), 10, 64)
	discriminator, _ := userData["discriminator"].(string)

	guildData, err := utils.HttpGet("https://discord.com/api/v6/users/@me/guilds", map[string]string{
		"Authorization": "Bearer " + response["access_token"].(string),
	}, true)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println(guildData["guilds"])

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
		log.Println(guild)
		var tmp database.GuildConfig
		r := database.Conn.Table("guild-configs").Where("\"GuildId\" = ?::bigint", guild["id"]).First(&tmp)
		if !errors.Is(r.Error, gorm.ErrRecordNotFound) {
			toSave := database.Guild{
				GuildID:    uint64(dId),
				Name:       guild["name"].(string),
				Icon:       icon,
				HasBot:     true,
				OwnerId:    -1,
				HasPremium: false,
			}
			tryCacheGuild(toSave)
			guildIDs = append(guildIDs, dId)
		}
		var pTmp database.Permissions
		r = database.Conn.Model(&database.Permissions{}).Where("guild = ? and user = ?", dId, userData["id"].(string)).Find(&pTmp)
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			database.Conn.Model(&database.Permissions{}).Create(&database.Permissions{
				Guild: dId,
				User:  userId,
				Perms: guild["permissions"].(float64),
			})
		} else {

			pTmp.Perms = guild["permissions"].(float64)
			database.Conn.Model(&database.Permissions{}).Where("guild = ? and user = ?", dId, userData["id"].(string)).Save(&pTmp)
		}
	}
	log.Printf("Caching guilds done for user \"%s\"", userData["username"].(string))

	if !initial {
		user.UserID = userId
		user.UserName = userData["username"].(string)
		user.Discriminator = discriminator
		user.AvatarID = userData["avatar"].(string)
		user.Guilds = utils.ArrayToPSQL(guildIDs)
		database.Conn.Model(&database.User{}).Where("user_id = ?", userId).Save(&user)
		return
	}
	userDb := database.User{
		UserID:        userId,
		UserName:      userData["username"].(string),
		Discriminator: discriminator,
		AvatarID:      userData["avatar"].(string),
		Guilds:        utils.ArrayToPSQL(guildIDs),
	}

	database.Conn.Create(&userDb)
}

func tryCacheGuild(guild database.Guild) {
	var tmp database.Guild
	r := database.Conn.Model(&database.Guild{}).First(&tmp, guild.GuildID)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		database.Conn.Model(&database.Guild{}).Create(&guild)
	} else {
		database.Conn.Model(&database.Guild{}).Where("gid = ?", guild.GuildID).Save(&guild)
	}
}
