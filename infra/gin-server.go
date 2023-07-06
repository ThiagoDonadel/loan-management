package infra

import (
	"log"
	"net/http"

	"github.com/ThiagoDonadel/loan-management/app/registry"
	"github.com/gin-gonic/gin"
)

var ginServer *gin.Engine

func StartGinServer() {
	ginServer = gin.Default()
	defineRoutes()

	log.Fatal(http.ListenAndServe(":8080", ginServer))
}

func defineRoutes() {
	routeGroup := ginServer.Group("/funding-calculator")
	registry.LoanController.SetupRoutes(routeGroup)
}
