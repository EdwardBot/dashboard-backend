package controllers

import (
	"github.com/edward-backend/controllers/auth"
	"github.com/gin-gonic/gin"
)

func InitAuth(router *gin.RouterGroup) {
	router.GET("/login", auth.HandleLogin)
	router.GET("/oauth", auth.HandleOAuth)
	router.GET("/user", auth.HandleUserInfo)
	router.Use(auth.HasAuth).POST(`/refresh`, auth.HandleRefresh)
}
