package app

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"users-api/controllers/users"
	"users-api/docs"
)

func mapUrls() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		eg := v1.Group("")
		{
			eg.POST("/users", users.Create)
			eg.GET("/users/:user_id", users.Get)
			eg.GET("/users/all", users.GetAll)
			eg.PUT("/users/:user_id", users.Update)   //Full update
			eg.PATCH("/users/:user_id", users.Update) //Partial Update
			eg.DELETE("/users/:user_id", users.Delete)
			eg.GET("/internal/users/search", users.Search)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
