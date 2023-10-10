package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ThiagoDonadel/loan-management/internal/security"
	"github.com/ThiagoDonadel/loan-management/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Base controller that holds the obligatory method that expose gin routes
type GinController interface {
	SetupRoutes(unauthorized, ownerAuthorized *gin.RouterGroup)
}

var ginServer *gin.Engine

func StartGinServer(controllers []GinController) {
	port := fmt.Sprintf(":%v", viper.Get("web.port"))
	ginServer = gin.Default()

	routeGroup := ginServer.Group("/funding-calculator")

	unauthorized := routeGroup.Group("")
	ownerAuthorized := routeGroup.Group("/:" + utils.OWNER_ID_PARAM_NAME)
	ownerAuthorized.Use(security.Authenticate())

	for _, controoler := range controllers {
		controoler.SetupRoutes(unauthorized, ownerAuthorized)
	}

	log.Fatal(http.ListenAndServe(port, ginServer))
}
