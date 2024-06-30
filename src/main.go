package main

import (
	"os"

	"github.com/Artexus/api-widyabhuvana/docs"
	"github.com/Artexus/api-widyabhuvana/src/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initSwagger(r *gin.Engine) {
	if os.Getenv("ENV") != "Production" {
		docs.SwaggerInfo.Title = "Swagger API Documentation"
		docs.SwaggerInfo.Description = "This is API documentation for api-widyabhuvana."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = os.Getenv("SERVER_HOST")
		docs.SwaggerInfo.Schemes = []string{"http"}

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}))

	initSwagger(r)
	route.InitRouter(r)

	r.Run(":" + os.Getenv("SERVER_PORT"))
}
