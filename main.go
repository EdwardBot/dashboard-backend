package main

import (
	"fmt"
	"github.com/edward-backend/controllers"
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-contrib/multitemplate"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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
	godotenv.Load()
	s = gocron.NewScheduler(time.UTC)
	router := gin.Default()
	err := database.Connect()
	if err != nil {
		panic(fmt.Sprintf("Error: %s", err.Error()))
	}
	utils.InitHttp()

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

	router.Run(":" + os.Getenv("PORT"))
}
