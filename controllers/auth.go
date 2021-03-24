package controllers

import (
	"github.com/edward-backend/controllers/auth"
	"github.com/gin-gonic/gin"
)

func InitAuth(router *gin.RouterGroup) {
	router.GET("/oauth", auth.HandleOAuth)
	router.GET("/refresh", auth.HandleRefresh)
	router.GET("/login", auth.HandleLogin)
}
