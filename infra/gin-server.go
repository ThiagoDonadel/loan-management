package infra

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ThiagoDonadel/loan-management/app/defaults"
	"github.com/ThiagoDonadel/loan-management/app/registry"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var ginServer *gin.Engine

func StartGinServer() {
	ginServer = gin.Default()
	defineRoutes()

	port := fmt.Sprintf(":%v", viper.Get("web.port"))

	log.Fatal(http.ListenAndServe(port, ginServer))
}

func defineRoutes() {
	routeGroup := ginServer.Group("/funding-calculator")

	unauthorized := routeGroup.Group("")
	ownerAuthorized := routeGroup.Group("/:" + defaults.OWNER_ID_PARAM_NAME)
	ownerAuthorized.Use(Authenticate())

	registry.LoanController.SetupRoutes(unauthorized, ownerAuthorized)
}
