package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		start := time.Now()

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()

		log.Printf("%s %s | %d | %v | %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			status,
			latency,
			ctx.ClientIP(),
		)
	}
}
