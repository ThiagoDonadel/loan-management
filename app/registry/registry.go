package registry

import (
	"github.com/ThiagoDonadel/loan-management/app/loan"
	"gorm.io/gorm"
)

var (
	loanRepository loan.Repository
)

var (
	loanService loan.Service
)

var (
	LoanController loan.Controller
)

func initializeRepositories(db *gorm.DB) {
	loanRepository = loan.NewRepository(db)
}

func initializeServices() {
	loanService = loan.NewService(loanRepository)
}

func initializeControllers() {
	LoanController = loan.NewController(loanService)
}

func Initialialize(db *gorm.DB) {
	initializeRepositories(db)
	initializeServices()
	initializeControllers()
}
