// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/FRahimov84/throttler/internal/usecase"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/FRahimov84/throttler/docs"
)

// NewRouter
// Swagger spec:
// @title       Throttler API
// @description Some desc
// @version     1.0
// @host        localhost:8080
// @BasePath	/api/v1
func NewRouter(handler *gin.Engine, t usecase.Throttler, l *zap.Logger) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	// Options

	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := handler.Group("/api/v1")
	{
		newThrottlerRoutes(v1, t, l)
	}
}
