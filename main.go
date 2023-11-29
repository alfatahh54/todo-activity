package main

import (
	"github.com/alfatahh54/todo-activity/routes"
	"github.com/alfatahh54/todo-activity/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	routes.Router(r)
	port := ":8080"
	getPort := utils.GoDotEnvVariable("PORT")
	if getPort != "" {
		port = ":" + getPort
	}
	r.Run(port)
}
