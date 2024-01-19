package services

import (
	loancalculator "github.com/ThiagoDonadel/loan-calculator"
	"github.com/ThiagoDonadel/loan-management/internal/metrics"
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/ThiagoDonadel/loan-management/internal/repository"
)

// LoanService is responsible to manage the bussines logic of loans
type LoanService interface {
	//Simulates a new Loan. This method will only calculate the valus and return
	Simulate(loan model.Loan) ([]model.LoanValue, error)
	//Calculate and creates a new loan
	Contract(loan model.Loan) (*model.Loan, error)
	//Find a loan by its ID
	Find(id, ownerId string) (*model.Loan, error)
	//Find all loans
	FindAll(ownerId string) ([]model.Loan, error)
}

type loanService struct {
	repo repository.LoanRepository
}

func NewLoanService(repository repository.LoanRepository) LoanService {
	return &loanService{
		repo: repository,
	}
}

func (s *loanService) Simulate(loan model.Loan) ([]model.LoanValue, error) {

	values, err := s.calculate(loan)

	if err != nil {
		return nil, err
	}

	metrics.IncrementSimulatedAmmountCounter(loan.LoanValue)

	return values, nil
}

func (s *loanService) Contract(loan model.Loan) (*model.Loan, error) {

	values, err := s.calculate(loan)

	if err != nil {
		return nil, err
	}

	for index := range values {
		values[index].OwnerId = loan.OwnerId
	}

	loan.Values = values
	s.repo.Create(&loan)

	return &loan, nil
}

func (s *loanService) calculate(loan model.Loan) ([]model.LoanValue, error) {

	params := loan.ToLoanCalcParameters()

	calculatedValues, err := loancalculator.Calculate(params)

	if err != nil {
		return nil, err
	}

	values := []model.LoanValue{}

	for _, calculatedValue := range calculatedValues {
		value := model.ConvertCalculatedValue(*calculatedValue)
		values = append(values, value)
	}

	return values, nil
}

func (s *loanService) Find(id, ownerId string) (*model.Loan, error) {

	foundLoan, err := s.repo.FindByID(id, ownerId)

	if err != nil {
		return nil, err
	}

	if foundLoan == nil {
		return nil, model.ErrLoanNotFound
	}

	return foundLoan, nil
}

func (s *loanService) FindAll(ownerId string) ([]model.Loan, error) {
	return s.repo.FindAll(ownerId)
}
