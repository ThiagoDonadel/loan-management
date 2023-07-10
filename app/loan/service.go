package loan

import (
	"errors"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
)

var (
	ErrLoanNotFound = errors.New("loan not found")
)

// Service is responsible to manage the bussines logic of loans
type Service interface {
	//Simulates a new Loan. This method will only calculate the valus and return
	Simulate(loan Loan) ([]LoanValue, error)
	//Calculate and creates a new loan
	Contract(loan Loan) (*Loan, error)
	//Find a loan by its ID
	Find(id string) (*Loan, error)
	//Find all loans
	FindAll() ([]Loan, error)
}

type service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return &service{
		repo: repository,
	}
}

func (s *service) Simulate(loan Loan) ([]LoanValue, error) {

	values, err := s.calculate(loan)

	if err != nil {
		return nil, err
	}

	return values, nil
}

func (s *service) Contract(loan Loan) (*Loan, error) {

	values, err := s.calculate(loan)

	if err != nil {
		return nil, err
	}

	loan.Values = values
	s.repo.Create(&loan)

	return &loan, nil
}

func (s *service) calculate(loan Loan) ([]LoanValue, error) {

	params := loan.toLoanCalcParameters()

	calculatedValues, err := loancalculator.Calculate(params)

	if err != nil {
		return nil, err
	}

	values := []LoanValue{}

	for _, calculatedValue := range calculatedValues {
		value := convertCalculatedValue(*calculatedValue)
		values = append(values, value)
	}

	return values, nil
}

func (s *service) Find(id string) (*Loan, error) {

	foundLoan, err := s.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	if foundLoan == nil {
		return nil, ErrLoanNotFound
	}

	return foundLoan, nil
}

func (s *service) FindAll() ([]Loan, error) {
	return s.repo.FindAll()
}
