package mocks

import (
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type LoanRepositoryMock struct {
	mock.Mock
}

func (m *LoanRepositoryMock) Create(newLoan *model.Loan) error {
	args := m.Called()
	randUUID := uuid.NewString()
	newLoan.Id = &randUUID
	return args.Error(0)
}

func (m *LoanRepositoryMock) FindByID(id, ownerId string) (*model.Loan, error) {
	args := m.Called()

	if args.Get(0) != nil {
		return args.Get(0).(*model.Loan), args.Error(1)
	} else {
		return nil, args.Error(1)
	}

}

func (m *LoanRepositoryMock) FindAll(ownerId string) ([]model.Loan, error) {
	args := m.Called()
	return args.Get(0).([]model.Loan), args.Error(1)
}
