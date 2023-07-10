package defaults

import "github.com/gin-gonic/gin"

//Base controller that holds the obligatory method that expose gin routes
type GinController interface {
	SetupRoutes(routerGroup *gin.RouterGroup)
}
