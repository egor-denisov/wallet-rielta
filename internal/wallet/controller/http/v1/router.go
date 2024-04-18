package v1

import (
	"log/slog"
	"net/http"

	_ "github.com/egor-denisov/wallet-rielta/docs" //nolint:blank-imports // for correct work swagger documentation
	"github.com/egor-denisov/wallet-rielta/internal/wallet/usecase"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger spec:
// @title       Wallet
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
// .
func NewRouter(handler *gin.Engine, l *slog.Logger, w usecase.Wallet) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/api/v1")
	{
		newWalletRoutes(h, w, l)
	}
}
