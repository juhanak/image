package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/juhanak/image/controllers"
)

func SetupRouters(app *gin.Engine) {
	api := app.Group("/api")
	api.GET("/images", controllers.GetImage)
	api.POST("/images", controllers.Post)
}
