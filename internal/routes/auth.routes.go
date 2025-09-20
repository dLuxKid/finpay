package router

import (
	"finpay/internal/handlers"
	"finpay/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, mongo *db.MongoInstance) {
	authHandler := handlers.NewAuthHandler(mongo)

	auth := r.Group("/auth")

	{
		auth.POST("/login", authHandler.Login)

		auth.POST("/register", authHandler.CreateUser)

		auth.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Logout endpoint"})
		})

		auth.POST("/refresh", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Refresh token endpoint"})
		})

		auth.POST("/forgot-password", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Forgot password endpoint"})
		})

		auth.POST("/reset-password", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "Reset password endpoint"})
		})
	}
}
