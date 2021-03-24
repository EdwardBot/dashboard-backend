package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	url2 "net/url"
	"os"
	"strings"
)

func HandleLogin(c *gin.Context) {
	url := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		os.Getenv("CLIENT_ID"), strings.ReplaceAll(url2.PathEscape(os.Getenv("BASE_URL")+"/v1/auth/oauth"), ":", "%3A"),
		"identify%20guilds")
	c.Redirect(308, url)
}
