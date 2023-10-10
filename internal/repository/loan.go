package repository

import (
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"gorm.io/gorm"
)

const CONTEXT_NAME = "loan-repo"

// LoanRepository is responsible to manage loan and loan valus in database
type LoanRepository interface {
	//Create a new loan in the database
	Create(newLoan *model.Loan) error
	//Find and retrieve a loan by its ID
	FindByID(id, ownerId string) (*model.Loan, error)
	//Find all loans
	FindAll(ownerId string) ([]model.Loan, error)
}

// Return a new instance of the repository
func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{
		db: db,
	}
}

type loanRepository struct {
	db *gorm.DB
}

func (r *loanRepository) Create(newLoan *model.Loan) error {

	if err := r.db.Create(newLoan).Error; err != nil {
		return err
	}

	return nil
}

func (r *loanRepository) FindByID(id, ownerId string) (*model.Loan, error) {

	foundLoan := &model.Loan{}

	preload := r.db.Preload("Values", func(db *gorm.DB) *gorm.DB {
		return db.Order("payment_date")
	})

	if err := preload.First(foundLoan, "id = ? and owner_id = ?", id, ownerId).Error; err != nil {
		return nil, err
	}

	return foundLoan, nil

}

func (r *loanRepository) FindAll(ownerId string) ([]model.Loan, error) {

	loans := []model.Loan{}

	if err := r.db.Where("owner_id = ?", ownerId).Find(&loans).Error; err != nil {
		return nil, err
	}

	return loans, nil
}
