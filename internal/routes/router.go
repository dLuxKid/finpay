package router

import (
	"finpay/internal/middlewares"
	"finpay/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routes(mongo *db.MongoInstance) *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.Logger())
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello from the server side")
	})

	AuthRoutes(r, mongo)
	UserRoutes(r)

	return r
}
