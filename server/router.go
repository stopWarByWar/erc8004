package server

import (
	"agent_identity/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitLogger(logger *logger.Logger) {
	_logger = logger
}

func Run(_cors []string, port string) {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     _cors,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Content-Disposition", "multipart/form-data"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("agent/identity/list", GetAgentCardListHandler)
	r.GET("agent/identity/filter/trustMode", GetAgentCardListByTrustModelHandler)
	r.GET("agent/identity/filter/Detail", GetAgentCardDetailHandler)

	r.Run(fmt.Sprintf(":%s", port))
}
