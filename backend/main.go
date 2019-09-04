package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zls/wiw-code-challenge/backend/model"
	"github.com/zls/wiw-code-challenge/backend/routes"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(model.NewDynamoDB())
	router.Use(model.NewMockSession())

	api := routes.Routes{}
	api.Register(router)

	router.Run(":8181")
}
