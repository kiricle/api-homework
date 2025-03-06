package router

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/kiricle/api-homework/docs"
	"github.com/kiricle/api-homework/internal/handlers"
	"github.com/swaggo/files" // swagger embed files
	"github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

func SetupRouter(bookHandler *handlers.BookHandler) *gin.Engine {
	r := gin.Default()

	r.Use(timeoutMiddleware(time.Minute * 2))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/book", bookHandler.GetBooks)
	r.POST("/book", bookHandler.CreateBook)
	r.GET("/book/:id", bookHandler.GetBook)

	return r
}

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel() // Ensure cancel is called to release resources

		// Replace the request's context with the new timeout context
		c.Request = c.Request.WithContext(ctx)

		// Run handler
		c.Next()

		// Check if context was canceled due to timeout
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
			c.Abort()
		}
	}
}
