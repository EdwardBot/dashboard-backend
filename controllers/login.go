package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//HandleLogin handles the Oauth2 login process
func HandleLogin(c *gin.Context) {
	authCode, hasCode := c.GetQuery("code")
	redirect, hasRedirect := c.GetQuery("redirect")
	if hasCode && hasRedirect {
		client := http.Client{}
		tokenRequest := url.Values{}
		tokenRequest.Set("client_id", "747157043466600477")
		tokenRequest.Set("client_secret", "8w6LQBDPD-Ypvrvr1meymB7gobg9OsDP")
		tokenRequest.Set("grant_type", "authorization_code")
		tokenRequest.Set("code", authCode)
		tokenRequest.Set("redirect_uri", redirect)
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
			c.JSON(401, map[string]interface{}{
				"status":  "Error",
				"error":   "API error!",
				"details": err,
			})
			return
		}
		log.Println(string(body))
		var response map[string]interface{}
		json.Unmarshal(body, &response)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, response)
		res.Body.Close()
	} else {
		c.JSON(401, map[string]interface{}{
			"status": "Error",
			"error":  "No code sent!",
		})
	}
}
