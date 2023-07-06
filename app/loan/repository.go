package loan

import "gorm.io/gorm"

const CONTEXT_NAME = "loan-repo"

type Repository interface {
	Create(newLoan *Loan) error
	FindByID(id string) (*Loan, error)
	FindAll() ([]Loan, error)
	PayInstallment(loanId string, installmentId uint64) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(newLoan *Loan) error {

	if err := r.db.Create(newLoan).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByID(id string) (*Loan, error) {

	foundLoan := &Loan{}

	preload := r.db.Preload("Values", func(db *gorm.DB) *gorm.DB {
		return db.Order("payment_date")
	})

	if err := preload.First(foundLoan, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return foundLoan, nil

}

func (r *repository) FindAll() ([]Loan, error) {

	loans := []Loan{}

	if err := r.db.Find(&loans).Error; err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *repository) PayInstallment(loanId string, installmentId uint64) error {

	if err := r.db.Model(&LoanValue{}).Where("id = ? AND loan_id = ?", installmentId, loanId).Update("paid", true).Error; err != nil {
		return err
	}

	return nil
}
