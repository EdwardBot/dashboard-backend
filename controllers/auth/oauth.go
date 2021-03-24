package auth

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
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

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			c.JSON(401, gin.H{
				"status":  "Error",
				"error":   "API error!",
				"details": err,
			})
			return
		}
		var response map[string]string
		json.Unmarshal(body, &response)
		if response["error_description"] == "Invalid \"code\" in request." {
			c.JSON(403, gin.H{
				"status": "error",
				"error":  "Code expired.",
			})
			return
		}

		req, _ := http.NewRequest("GET", "https://discord.com/api/v6/users/@me", strings.NewReader(""))

		req.Header.Add("Authorization", "Bearer "+response["access_token"])

		res, _ = client.Do(req)
		body, _ = ioutil.ReadAll(res.Body)
		var userData map[string]string
		json.Unmarshal(body, &userData)
		log.Println(string(body))

		token := jwt.New(jwt.GetSigningMethod("HS256"))

		token.Claims = jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Subject:   userData["id"],
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
