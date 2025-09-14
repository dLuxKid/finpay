package router

import "github.com/gin-gonic/gin"

func Routes() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello from the server side")
	})

	UserRoutes(r)

	return r
}
