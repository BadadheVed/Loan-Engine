package router

import (
	"github.com/BadadheVed/clickpe/controllers"
	"github.com/gin-gonic/gin"
)

func apiRouter(r *gin.Engine) {

	api := r.Group("/api")
	api.GET("/health", controllers.Health)
	api.POST("/uploadcsv", controllers.UploadCSVUsers)

}
