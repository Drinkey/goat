package routers

import (
	_ "github.com/Drinkey/goat/docs"
	v1 "github.com/Drinkey/goat/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const currentAPIVersion = "/api/v1"

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV1 := r.Group(currentAPIVersion)
	{
		apiV1.GET("/ping", v1.Ping)
		apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		apiV1.GET("/cron", v1.ListCronTasks)
		apiV1.POST("/cron/:id", v1.RunOneTask)
		apiV1.GET("/cron/:id", v1.GetOneTask)
	}
	return r
}
