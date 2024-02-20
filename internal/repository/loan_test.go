package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/ThiagoDonadel/loan-management/internal/repository"
	"github.com/ThiagoDonadel/loan-management/test/containers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoanRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LoanRepositoryTestSuite))
}

//----- TEST

type LoanRepositoryTestSuite struct {
	suite.Suite
}

func (s *LoanRepositoryTestSuite) TestA() {
	pgContainer, err := containers.CreatePostgressTestContainer("loan-test", "admin", "123456")

	if err != nil {
		assert.Fail(s.T(), err.Error())
	}

	defer pgContainer.Terminate(context.Background())

	port, _ := pgContainer.MappedPort(context.Background(), "5432")

	fmt.Println(port)

	DSN := fmt.Sprintf(
		"user=%v password=%v dbname=%v port=%v sslmode=disable",
		"admin",
		"123456",
		"loan-test",
		port.Int(),
	)

	fmt.Println(pgContainer.ConnectionString(context.Background(), "sslmode=disable"))

	config := postgres.Config{
		DSN:                  DSN,
		PreferSimpleProtocol: true,
	}

	db, err := gorm.Open(postgres.New(config), &gorm.Config{})

	if err != nil {
		assert.Fail(s.T(), err.Error())
	}

	err = db.AutoMigrate(&model.Loan{}, &model.LoanValue{})

	if err != nil {
		assert.Fail(s.T(), err.Error())
	}

	loan := model.Loan{
		Method:         int(loancalculator.CONSTANT_AMORTIZATION),
		LoanValue:      5000,
		Rate:           0.5,
		Term:           12,
		RateBaseMonths: 1,
		StartDate:      model.LocalDate(time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local)),
		OwnerId:        "aaaaaaa", //uuid.NewString(),
		Values: []model.LoanValue{
			{Number: 1,
				PaymentDate:  model.LocalDate(time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local)),
				Installment:  10.5,
				Amortization: 10.5,
				Interest:     10.5,
				Balance:      90.5,
			},
		},
	}

	repo := repository.NewLoanRepository(db)
	err = repo.Create(&loan)

	if err != nil {
		assert.Fail(s.T(), err.Error())
	}

	assert.NotNil(s.T(), loan.Id)
}
