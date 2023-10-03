package loan

import (
	"errors"
	"net/http"

	"github.com/ThiagoDonadel/loan-management/app/defaults"
	"github.com/gin-gonic/gin"
)

// Controller is responsible to manage the web logic of loans
type Controller interface {
	//Simulates a new Loan. This method will only calculate the valus and return
	Simulate(context *gin.Context)
	//Calculate and creates a new loan
	Contract(context *gin.Context)
	//Find a loan by its ID
	Find(context *gin.Context)
	//Find all loans
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

	if err := context.ShouldBindJSON(params); err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	values, _ := c.service.Simulate(*params)

	context.JSON(http.StatusOK, values)
}

func (c *controller) Contract(context *gin.Context) {
	params := &Loan{}

	if err := context.ShouldBindJSON(params); err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	params.OwnerId = context.Param(defaults.OWNER_ID_PARAM_NAME)
	values, _ := c.service.Contract(*params)

	context.JSON(http.StatusOK, values)
}

func (c *controller) Find(context *gin.Context) {

	loanFound, err := c.service.Find(context.Param(defaults.RESOURCE_ID_PARAM_NAME), context.Param(defaults.OWNER_ID_PARAM_NAME))

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

	loans, err := c.service.FindAll(context.Param(defaults.OWNER_ID_PARAM_NAME))

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	context.JSON(http.StatusOK, loans)
}

func (c *controller) SetupRoutes(unauthorized, ownerAuthorized *gin.RouterGroup) {
	ownerAuthorized.POST("/simulate/", c.Simulate)
	ownerAuthorized.POST("/contract/", c.Contract)
	ownerAuthorized.GET("/find/:"+defaults.RESOURCE_ID_PARAM_NAME+"/", c.Find)
	ownerAuthorized.GET("/find-all/", c.FIndAll)
}
