package registry

import (
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/ThiagoDonadel/loan-management/internal/repository"
	"github.com/ThiagoDonadel/loan-management/internal/services"
	"github.com/ThiagoDonadel/loan-management/internal/web"
	"gorm.io/gorm"
)

var (
	loanRepository repository.LoanRepository
)

var (
	loanService services.LoanService
)

var (
	LoanController web.LoanController
)

func initializeRepositories(db *gorm.DB) {
	loanRepository = repository.NewLoanRepository(db)
}

func initializeServices() {
	loanService = services.NewLoanService(loanRepository)
}

func initializeControllers() {
	LoanController = web.NewLoanController(loanService)
}

func Initialialize(db *gorm.DB) {
	initializeRepositories(db)
	initializeServices()
	initializeControllers()
}

func GetControllers() []model.BaseController {
	return []model.BaseController{LoanController}
}
