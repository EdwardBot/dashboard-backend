package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/edward-backend/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

var (
	s *gocron.Scheduler
)

func main() {
	log.Println(`Starting...`)
	s = gocron.NewScheduler(time.UTC)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
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

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, map[string]interface{}{
			"status": "error",
			"error":  "Route not found!",
		})
	})

	r.GET("/auth/login", controllers.HandleLogin)

	if os.Getenv("PORT") != "" {
		gin.SetMode(gin.ReleaseMode)
		r.Run(":" + os.Getenv("PORT"))
		return
	} else {
		gin.SetMode(gin.ReleaseMode)
		r.Run(":3000")
	}
}
