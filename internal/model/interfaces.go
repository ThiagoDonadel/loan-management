package model

import "github.com/gin-gonic/gin"

// Base controller that holds the obligatory method that expose gin routes
type BaseController interface {
	SetupRoutes(unauthorized, ownerAuthorized *gin.RouterGroup)
}
