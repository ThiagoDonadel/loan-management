package defaults

import "github.com/gin-gonic/gin"

type GinController interface {
	SetupRoutes(routerGroup *gin.RouterGroup)
}
