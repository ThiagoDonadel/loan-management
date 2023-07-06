package loan

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ThiagoDonadel/loan-management/app/defaults"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Simulate(context *gin.Context)
	Contract(context *gin.Context)
	Pay(context *gin.Context)
	Find(context *gin.Context)
	FIndAll(context *gin.Context)
	defaults.GinController
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

type controller struct {
	service Service
}

func (c *controller) Simulate(context *gin.Context) {

	params := &Loan{}
	err := context.ShouldBindJSON(params)
	fmt.Println(err)
	values, _ := c.service.Simulate(*params)

	context.JSON(http.StatusOK, values)
}

func (c *controller) Contract(context *gin.Context) {
	params := &Loan{}
	err := context.ShouldBindJSON(params)
	fmt.Println(err)
	values, _ := c.service.Contract(*params)

	context.JSON(http.StatusOK, values)
}

func (c *controller) Pay(context *gin.Context) {

	installmentId, _ := strconv.ParseUint(context.Param("installmentId"), 10, 64)

	c.service.Pay(context.Param("id"), installmentId)

}

func (c *controller) Find(context *gin.Context) {

	loanFound, err := c.service.Find(context.Param("id"))

	if err != nil {
		if errors.Is(err, ErrLoanNotFound) {
			context.JSON(http.StatusNotFound, err.Error())
		} else {
			context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	context.JSON(http.StatusOK, loanFound)

}

func (c *controller) FIndAll(context *gin.Context) {

	loans, err := c.service.FindAll()

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	context.JSON(http.StatusOK, loans)
}

func (c *controller) SetupRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/simulate/", c.Simulate)
	routerGroup.POST("/contract/", c.Contract)
	routerGroup.POST("/pay/:id/:installmentId", c.Pay)
	routerGroup.GET("/find/:id/", c.Find)
	routerGroup.GET("/find-all/", c.FIndAll)

}
