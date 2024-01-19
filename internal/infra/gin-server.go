package infra

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ThiagoDonadel/loan-management/internal/metrics"
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/ThiagoDonadel/loan-management/internal/security"
	"github.com/ThiagoDonadel/loan-management/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var ginServer *gin.Engine

func StartGinServer(controllers []model.BaseController) {
	port := fmt.Sprintf(":%v", viper.Get("web.port"))
	ginServer = gin.Default()

	ginServer.GET("/metrics", gin.WrapH(promhttp.Handler()))

	routeGroup := ginServer.Group("/funding-calculator")
	routeGroup.Use(metrics.PrometheusMetricsHandler())

	unauthorized := routeGroup.Group("")
	ownerAuthorized := routeGroup.Group("/:" + utils.OWNER_ID_PARAM_NAME)
	ownerAuthorized.Use(security.Authenticate())

	for _, controoler := range controllers {
		controoler.SetupRoutes(unauthorized, ownerAuthorized)
	}

	log.Fatal(http.ListenAndServe(port, ginServer))
}
