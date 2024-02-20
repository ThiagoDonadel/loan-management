package services_test

import (
	"errors"
	"testing"
	"time"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/ThiagoDonadel/loan-management/internal/services"
	"github.com/ThiagoDonadel/loan-management/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestLoanServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LoanServiceTestSuite))
}

//----- TEST

type LoanServiceTestSuite struct {
	suite.Suite
	loanId        string
	ownerId       string
	validLoan     model.Loan
	invalidLoan   model.Loan
	databaseLoan  model.Loan
	databaseError error
}

func (s *LoanServiceTestSuite) SetupSuite() {

	s.loanId = uuid.NewString()
	s.ownerId = uuid.NewString()

	s.validLoan = model.Loan{
		Method:         int(loancalculator.CONSTANT_AMORTIZATION),
		LoanValue:      5000,
		Rate:           0.5,
		Term:           12,
		RateBaseMonths: 1,
		StartDate:      model.LocalDate(time.Date(2023, 2, 1, 0, 0, 0, 0, time.Local)),
		OwnerId:        s.ownerId,
	}

	s.invalidLoan = model.Loan{
		Method:         10,
		LoanValue:      5000,
		Rate:           0.5,
		Term:           12,
		RateBaseMonths: 1,
		StartDate:      model.LocalDate(time.Date(2023, 2, 1, 0, 0, 0, 0, time.Local)),
		OwnerId:        s.ownerId,
	}

	s.databaseLoan = model.Loan{
		Id:             &s.loanId,
		Method:         int(loancalculator.CONSTANT_AMORTIZATION),
		LoanValue:      5000,
		Rate:           0.5,
		Term:           12,
		RateBaseMonths: 1,
		StartDate:      model.LocalDate(time.Date(2023, 2, 1, 0, 0, 0, 0, time.Local)),
		OwnerId:        s.ownerId,
	}

	s.databaseError = errors.New("Database Error")

}

func (s *LoanServiceTestSuite) Test_When_Simulate_Success() {

	expectedLines := 13
	service := services.NewLoanService(&mocks.LoanRepositoryMock{})

	result, err := service.Simulate(s.validLoan)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedLines, len(result))
	for index, value := range result {
		assert.Equal(s.T(), index, value.Number)
	}
}

func (s *LoanServiceTestSuite) Test_When_Simulate_ValidationError() {

	service := services.NewLoanService(&mocks.LoanRepositoryMock{})

	_, err := service.Simulate(s.invalidLoan)

	assert.Error(s.T(), err)
}

func (s *LoanServiceTestSuite) Test_When_Contract_Success() {

	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("Create").Return(nil)
	service := services.NewLoanService(repoMock)

	result, err := service.Contract(s.validLoan)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), result.Id)
	for index, value := range result.Values {
		assert.Equal(s.T(), index, value.Number)
		assert.NotEmpty(s.T(), value.OwnerId)
	}
}

func (s *LoanServiceTestSuite) Test_When_Contract_ValidationError() {

	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("Create").Return(nil)
	service := services.NewLoanService(repoMock)

	_, err := service.Contract(s.invalidLoan)

	assert.Error(s.T(), err)
}

func (s *LoanServiceTestSuite) Test_When_Contract_DatabaseError() {

	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("Create").Return(s.databaseError)
	service := services.NewLoanService(repoMock)

	_, err := service.Contract(s.validLoan)

	assert.ErrorIs(s.T(), err, s.databaseError)
}

func (s *LoanServiceTestSuite) Test_When_Find_Success() {

	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("FindByID").Return(&s.databaseLoan, nil)
	service := services.NewLoanService(repoMock)

	result, err := service.Find(s.loanId, s.ownerId)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), &s.databaseLoan, result)
}

func (s *LoanServiceTestSuite) Test_When_Find_DatabaseError() {

	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("FindByID").Return(nil, s.databaseError)
	service := services.NewLoanService(repoMock)

	_, err := service.Find(s.loanId, s.ownerId)

	assert.ErrorIs(s.T(), err, s.databaseError)
}

func (s *LoanServiceTestSuite) Test_When_Find_NotFoundError() {

	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("FindByID").Return(nil, nil)
	service := services.NewLoanService(repoMock)

	_, err := service.Find(s.loanId, s.ownerId)

	assert.ErrorIs(s.T(), err, model.ErrLoanNotFound)
}

func (s *LoanServiceTestSuite) Test_When_FindAll_Success() {

	expected := []model.Loan{s.databaseLoan}
	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("FindAll").Return(expected, nil)
	service := services.NewLoanService(repoMock)

	result, err := service.FindAll(s.loanId)

	assert.Nil(s.T(), err)
	assert.ElementsMatch(s.T(), expected, result)
}

func (s *LoanServiceTestSuite) Test_When_FindAll_DatabaseError() {

	repoMock := &mocks.LoanRepositoryMock{}
	repoMock.On("FindAll").Return([]model.Loan{}, s.databaseError)
	service := services.NewLoanService(repoMock)

	_, err := service.FindAll(s.ownerId)

	assert.ErrorIs(s.T(), err, s.databaseError)
}
