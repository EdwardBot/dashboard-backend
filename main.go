package main

import (
	"encoding/json"
	"github.com/edward-backend/controllers"
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-contrib/multitemplate"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

var (
	s *gocron.Scheduler
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("redirect", "templates/oauth.html")
	return r
}

func main() {
	log.Println(`Starting...`)
	if os.Getenv("PROD") == "" {
		godotenv.Load()
	}
	s = gocron.NewScheduler(time.UTC)
	router := gin.Default()
	database.Connect()

	database.Init()

	router.HTMLRender = createMyRender()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	utils.InitDiscord()

	router.Use(func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Content-Type"), "application/json") {
			c.Next()
		}
		rawBody, _ := io.ReadAll(c.Request.Body)
		var body map[string]interface{}
		err := json.Unmarshal(rawBody, &body)
		if err != nil {
			c.Next()
		}
		c.Set("body", body)
		c.Next()
	})

	r := router.Group("/v1")

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, map[string]interface{}{
			"status": "error",
			"error":  "Route not found!",
		})
	})

	controllers.InitAuth(r.Group("/auth"))
	controllers.InitGuilds(r.Group("/guild"))

	router.Run(":" + os.Getenv("PORT"))
}
