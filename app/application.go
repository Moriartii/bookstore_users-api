package app

import (
	"github.com/Moriartii/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()

	logger.Info("we begin start appication")
	router.Run(":8081")
}
