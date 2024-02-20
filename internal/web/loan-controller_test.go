package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/ThiagoDonadel/loan-management/internal/repository"
	"github.com/ThiagoDonadel/loan-management/internal/services"
	"github.com/ThiagoDonadel/loan-management/internal/utils"
	"github.com/ThiagoDonadel/loan-management/internal/web"
	"github.com/ThiagoDonadel/loan-management/test/containers"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LoanControllerTestSuite struct {
	suite.Suite
	pgContainer     *postgresContainer.PostgresContainer
	gormDB          *gorm.DB
	httpAuthRoute   string
	httpUnauthRoute string
	loanService     services.LoanService
	ownerId         string
	loanId          string
	loan            model.Loan
}

func TestLoanController(t *testing.T) {
	suite.Run(t, &LoanControllerTestSuite{})
}

func (s *LoanControllerTestSuite) SetupSuite() {

	var err error
	dbUser, dbPass, dbName, dbPort := "admin", "123456", "loan-test", "5432"

	s.pgContainer, err = containers.CreatePostgressTestContainer(dbName, dbUser, dbPass)

	if err != nil {
		assert.Fail(s.T(), err.Error())
	}

	port, _ := s.pgContainer.MappedPort(context.Background(), nat.Port(dbPort))

	DSN := fmt.Sprintf(
		"user=%v password=%v dbname=%v port=%v sslmode=disable",
		dbUser,
		dbPass,
		dbName,
		port.Int(),
	)

	config := postgres.Config{
		DSN:                  DSN,
		PreferSimpleProtocol: true,
	}

	s.gormDB, err = gorm.Open(postgres.New(config), &gorm.Config{})

	if err != nil {
		assert.Fail(s.T(), err.Error())
	}

	err = s.gormDB.AutoMigrate(&model.Loan{}, &model.LoanValue{})

	if err != nil {
		assert.Fail(s.T(), err.Error())
	}

	s.ownerId = uuid.NewString()

	s.loan = model.Loan{
		Method:         int(loancalculator.CONSTANT_AMORTIZATION),
		LoanValue:      5000,
		Rate:           0.5,
		Term:           12,
		RateBaseMonths: 1,
		StartDate:      model.LocalDate(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.Local)),
		OwnerId:        s.ownerId,
		Values: []model.LoanValue{
			{Number: 1,
				PaymentDate:  model.LocalDate(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.Local)),
				Installment:  10.5,
				Amortization: 10.5,
				Interest:     10.5,
				Balance:      90.5,
			},
		},
	}

	if err := s.gormDB.Create(&s.loan).Error; err != nil {
		assert.Fail(s.T(), err.Error())
	}

	s.loanId = *s.loan.Id
	s.loanService = services.NewLoanService(repository.NewLoanRepository(s.gormDB))

	s.httpUnauthRoute = "/testNoAuth"
	s.httpAuthRoute = "/testAuth/"

}

func (s *LoanControllerTestSuite) TearDownSuite() {
	db, _ := s.gormDB.DB()
	db.Close()
	s.pgContainer.Terminate(context.Background())
}

func (s *LoanControllerTestSuite) Test_When_Find() {
	var result model.Loan
	w := httptest.NewRecorder()
	controller := web.NewLoanController(s.loanService)
	server := gin.Default()
	controller.SetupRoutes(server.Group(s.httpUnauthRoute), server.Group(s.httpAuthRoute+":"+utils.OWNER_ID_PARAM_NAME))

	path := s.httpAuthRoute + s.ownerId + "/find/" + s.loanId
	request, _ := http.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Content-Type", "application/json")

	server.ServeHTTP(w, request)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &result)

	assert.Equal(s.T(), s.loanId, result.Id)
}

func (s *LoanControllerTestSuite) Test_When_Find_All() {
	var result []model.Loan
	controller := web.NewLoanController(s.loanService)
	w := httptest.NewRecorder()
	server := gin.Default()
	controller.SetupRoutes(server.Group(s.httpUnauthRoute), server.Group(s.httpAuthRoute+":"+utils.OWNER_ID_PARAM_NAME))

	path := s.httpAuthRoute + s.ownerId + "/find-all"
	request, _ := http.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Content-Type", "application/json")

	server.ServeHTTP(w, request)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &result)

	assert.Greater(s.T(), len(result), 0)
}

func (s *LoanControllerTestSuite) Test_When_Contract() {
	var result model.Loan
	controller := web.NewLoanController(s.loanService)
	w := httptest.NewRecorder()
	server := gin.Default()
	controller.SetupRoutes(server.Group(s.httpUnauthRoute), server.Group(s.httpAuthRoute+":"+utils.OWNER_ID_PARAM_NAME))

	toSave := s.loan
	toSave.Id = nil
	toSave.Values = nil
	reqBody, _ := json.Marshal(toSave)
	path := s.httpAuthRoute + s.ownerId + "/contract"
	request, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	server.ServeHTTP(w, request)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &result)

	assert.NotEmpty(s.T(), result.Id)
}

func (s *LoanControllerTestSuite) Test_When_Simulate() {
	var result []model.LoanValue
	controller := web.NewLoanController(s.loanService)
	w := httptest.NewRecorder()
	server := gin.Default()
	controller.SetupRoutes(server.Group(s.httpUnauthRoute), server.Group(s.httpAuthRoute+":"+utils.OWNER_ID_PARAM_NAME))

	toSimulate := s.loan
	toSimulate.Id = nil
	toSimulate.Values = nil
	reqBody, _ := json.Marshal(toSimulate)
	path := s.httpUnauthRoute + "/simulate"
	request, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	server.ServeHTTP(w, request)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	body, _ := io.ReadAll(w.Body)
	json.Unmarshal(body, &result)

	assert.NotEmpty(s.T(), result)
}
