package routers

import (
	_ "github.com/Drinkey/goat/docs"
	v1 "github.com/Drinkey/goat/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const currentAPIVersion = "/api/v1"

// CORS allows cross origin access
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), CORS())

	apiV1 := r.Group(currentAPIVersion)
	{
		apiV1.GET("/ping", v1.Ping)
		apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		apiV1.GET("/cron/", v1.ListCronTasks)
		apiV1.POST("/cron/:id", v1.RunOneTask)
		apiV1.GET("/cron/:id", v1.GetOneTask)
	}
	return r
}
