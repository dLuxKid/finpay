package middlewares

import (
	"finpay/internal/models"
	"slices"

	"github.com/gin-gonic/gin"
)

func RequiredRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userObj, exists := ctx.Get("user")

		if !exists {
			ctx.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized", "error": "User not found in context", "success": false})
			return
		}

		user := userObj.(*models.User)

		if slices.Contains(roles, user.Role) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(403, gin.H{"message": "Forbidden", "error": "User does not have the required role", "success": false})
	}
}
