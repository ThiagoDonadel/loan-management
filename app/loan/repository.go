package loan

import "gorm.io/gorm"

const CONTEXT_NAME = "loan-repo"

//Repository is responsible to manage loan and loan valus in database
type Repository interface {
	//Create a new loan in the database
	Create(newLoan *Loan) error
	//Find and retrieve a loan by its ID
	FindByID(id string) (*Loan, error)
	//Find all loans
	FindAll() ([]Loan, error)
}

//Return a new instance of the repository
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
