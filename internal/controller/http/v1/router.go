package v1

import (
	"go.uber.org/zap"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/FRahimov84/throttler/internal/usecase"
)

// NewRouter
// Swagger spec:
// @title       Throttler API
// @description Some desc
// @version     1.0
// @host        localhost:8080
// @BasePath	/api/v1
func NewRouter(handler *gin.Engine, t usecase.Throttler, l *zap.Logger) {
	// Options
	handler.Use(ginzap.Ginzap(l, time.RFC3339, true))
	handler.Use(ginzap.RecoveryWithZap(l, true))

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
