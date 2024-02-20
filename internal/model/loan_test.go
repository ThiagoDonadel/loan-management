package model_test

import (
	"testing"
	"time"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
	"github.com/ThiagoDonadel/loan-management/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_Convert_To_Calculation_Params(t *testing.T) {

	//preparation
	loan := model.Loan{
		Method:         int(loancalculator.CONSTANT_AMORTIZATION),
		LoanValue:      5000,
		Rate:           0.5,
		Term:           12,
		RateBaseMonths: 1,
		StartDate:      model.LocalDate(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.Local)),
	}

	//test
	result := loan.ToLoanCalcParameters()

	//asserts
	assert.Equal(t, loan.Method, int(result.Method))
	assert.Equal(t, loan.LoanValue, result.InitialValue)
	assert.Equal(t, loan.Rate, result.Rate)
	assert.Equal(t, loan.Term, result.Term)
	assert.Equal(t, loan.RateBaseMonths, int(result.RateBaseMonths))
	assert.Equal(t, time.Time(loan.StartDate), result.BaseDate)

}

func Test_Convert_Calculation_Value(t *testing.T) {

	//preparation
	value := loancalculator.Value{
		Number:       1,
		PaymentDate:  time.Date(2024, time.February, 1, 0, 0, 0, 0, time.Local),
		Installment:  10.5,
		Amortization: 10.5,
		Interest:     10.5,
		Balance:      90.5,
	}

	//test
	result := model.ConvertCalculatedValue(value)

	//asserts
	assert.Equal(t, value.Number, result.Number)
	assert.Equal(t, value.PaymentDate, time.Time(result.PaymentDate))
	assert.Equal(t, value.Installment, result.Installment)
	assert.Equal(t, value.Amortization, result.Amortization)
	assert.Equal(t, value.Interest, result.Interest)
	assert.Equal(t, value.Balance, result.Balance)

}
