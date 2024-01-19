package web

import (
	"errors"
	"net/http"

	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/ThiagoDonadel/loan-management/internal/services"
	"github.com/ThiagoDonadel/loan-management/internal/utils"
	"github.com/gin-gonic/gin"
)

// LoanController is responsible to manage the web logic of loans
type LoanController interface {
	//Simulates a new Loan. This method will only calculate the valus and return
	Simulate(context *gin.Context)
	//Calculate and creates a new loan
	Contract(context *gin.Context)
	//Find a loan by its ID
	Find(context *gin.Context)
	//Find all loans
	FIndAll(context *gin.Context)
	model.BaseController
}

func NewLoanController(service services.LoanService) LoanController {
	return &loanController{service: service}
}

type loanController struct {
	service services.LoanService
}

func (c *loanController) Simulate(context *gin.Context) {

	params := &model.Loan{}

	if err := context.ShouldBindJSON(params); err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	values, _ := c.service.Simulate(*params)

	context.JSON(http.StatusOK, values)
}

func (c *loanController) Contract(context *gin.Context) {
	params := &model.Loan{}

	if err := context.ShouldBindJSON(params); err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	params.OwnerId = context.Param(utils.OWNER_ID_PARAM_NAME)
	values, _ := c.service.Contract(*params)

	context.JSON(http.StatusOK, values)
}

func (c *loanController) Find(context *gin.Context) {

	loanFound, err := c.service.Find(context.Param(utils.RESOURCE_ID_PARAM_NAME), context.Param(utils.OWNER_ID_PARAM_NAME))

	if err != nil {
		if errors.Is(err, model.ErrLoanNotFound) {
			context.JSON(http.StatusNotFound, err.Error())
		} else {
			context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	context.JSON(http.StatusOK, loanFound)

}

func (c *loanController) FIndAll(context *gin.Context) {

	loans, err := c.service.FindAll(context.Param(utils.OWNER_ID_PARAM_NAME))

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
	}

	context.JSON(http.StatusOK, loans)
}

func (c *loanController) SetupRoutes(unauthorized, ownerAuthorized *gin.RouterGroup) {
	unauthorized.POST("/simulate/", c.Simulate)
	ownerAuthorized.POST("/contract/", c.Contract)
	ownerAuthorized.GET("/find/:"+utils.RESOURCE_ID_PARAM_NAME+"/", c.Find)
	ownerAuthorized.GET("/find-all/", c.FIndAll)
}
