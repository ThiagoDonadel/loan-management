package loan

import (
	"errors"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
)

var (
	ErrLoanNotFound = errors.New("loan not found")
)

type Service interface {
	Simulate(loan Loan) ([]LoanValue, error)
	Contract(loan Loan) (*Loan, error)
	Find(id string) (*Loan, error)
	FindAll() ([]Loan, error)
	Pay(id string, installmentId uint64) error
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
		if value.Number == 0 {
			tBool := true
			value.Paid = &tBool
		}
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

func (s *service) Pay(id string, installmentId uint64) error {
	return s.repo.PayInstallment(id, installmentId)
}
