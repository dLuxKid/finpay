package handlers

import (
	"finpay/internal/config"
	"finpay/internal/models"
	"finpay/pkg/db"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct to hold the mongo instance, constructor function to create a new AuthHandler with the mongo instance
type AuthHandler struct {
	// this makes sure every method in AuthHandler has access to the database
	mongo *db.MongoInstance
}

// constructor function
func NewAuthHandler(mongo *db.MongoInstance) *AuthHandler {
	// return a pointer to AuthHandler with the mongo instance, makes a new AuthHandler with the mongo instance
	return &AuthHandler{mongo: mongo}
}

// defines a method on AuthHandler, method to create a new user, receiver is a pointer to AuthHandler
func (h *AuthHandler) CreateUser(ctx *gin.Context) {
	var newUser models.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false, "message": "Invalid request body"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password", "error": err.Error()})
		return
	}

	newUser.Password = string(hashedPassword)
	newUser.ID = primitive.NewObjectID()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	res, err := h.mongo.Users.InsertOne(ctx, newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error(), "success": false})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   newUser.Username,
		"email":      newUser.Email,
		"id":         newUser.ID.Hex(),
		"expired_at": jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		"issued_at":  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(config.Load().JWTSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token", "error": err.Error(), "success": false})
		return
	}

	ctx.SetCookie("finpay_jwt_token", tokenString, 3600*72, "/", "", false, true)

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": newUser, "token": tokenString, "success": true, "id": res.InsertedID})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false, "message": "Invalid request body"})
		return
	}

	filter := bson.M{"email": loginData.Email}
	var user models.User
	err := h.mongo.Users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid email or password", "error": err.Error()})
		return
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if hashErr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "success": false, "error": hashErr.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   user.Username,
		"email":      user.Email,
		"id":         user.ID.Hex(),
		"expired_at": jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		"issued_at":  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(config.Load().JWTSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token", "error": err.Error(), "success": false})
		return
	}

	ctx.SetCookie("finpay_jwt_token", tokenString, 3600*72, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "token": tokenString, "success": true})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("finpay_jwt_token", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful", "success": true})
}

// func (h * AuthHandler) RefreshToken(ctx *gin.Context){}
