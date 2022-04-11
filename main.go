package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/juhanak/image/routers"
)

func main() {
	app := gin.Default()
	routers.SetupRouters(app)
	err := app.Run(":8080")
	if err != nil {
		fmt.Errorf("failed to setup server: %v", err)
	}
}
