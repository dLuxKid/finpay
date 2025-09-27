package middlewares

import (
	"finpay/internal/models"
	"finpay/pkg/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthMiddleware struct {
	DB        *db.MongoInstance
	JWTSecret []byte
}

func NewAuthMiddleware(db *db.MongoInstance, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		DB:        db,
		JWTSecret: []byte(jwtSecret),
	}
}

func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["id"].(string)

			var user models.User
			err := a.DB.Users.FindOne(c, bson.M{"_id": userID}).Decode(&user)
			if err == nil {
				c.Set("user", &user)
			}
		}

		c.Next()
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(("Authorization"))
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": "Missing Authorization header", "success": false})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader == tokenString {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": "Invalid Authorization header format", "success": false})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return a.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": "Invalid token", "success": false})

			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := claims["id"].(string)
			var user models.User

			if err := a.DB.Users.FindOne(ctx, bson.M{"_id": userId}).Decode(&user); err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": "User not found", "success": false})
				return
			}
			ctx.Set("user", user)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": "Invalid token claims", "success": false})
			return
		}

		ctx.Next()
	}
}
