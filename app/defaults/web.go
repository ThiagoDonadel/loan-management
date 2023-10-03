package defaults

import "github.com/gin-gonic/gin"

//Base controller that holds the obligatory method that expose gin routes
type GinController interface {
	SetupRoutes(unauthorized, ownerAuthorized *gin.RouterGroup)
}

//Path param names definitions
const RESOURCE_ID_PARAM_NAME = "id"
const OWNER_ID_PARAM_NAME = "ownerid"
