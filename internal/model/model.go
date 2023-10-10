package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
)

// Type used customize the time format
type LocalDate time.Time

func (l *LocalDate) MarshalJSON() ([]byte, error) {
	formatedDate := time.Time(*l).Format(time.DateOnly)
	return json.Marshal(formatedDate)
}

func (l *LocalDate) UnmarshalJSON(bytes []byte) error {

	dateString := strings.Trim(string(bytes), "\"")
	date, err := time.Parse(time.DateOnly, dateString)

	if err != nil {
		return ErrInvalidDateFormat
	}

	*l = LocalDate(date)

	return nil
}

func (l *LocalDate) Scan(value interface{}) error {
	date, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	*l = LocalDate(date)

	return nil
}

func (l LocalDate) Value() (driver.Value, error) {
	return driver.Value(time.Time(l)), nil
}

// Strusct that holds the loan data
type Loan struct {
	Id             *string     `json:"id" gorm:"column:id;primarykey;default:gen_random_uuid()"`
	OwnerId        string      `json:"-" gorm:"column:owner_id"`
	Method         int         `json:"method" gorm:"column:method"`
	LoanValue      float64     `json:"value" gorm:"column:loan_value"`
	Rate           float64     `json:"rate" gorm:"column:rate"`
	RateBaseMonths int         `json:"rate_base_months" gorm:"column:rate_base_months"`
	Term           int         `json:"term" gorm:"column:term"`
	StartDate      LocalDate   `json:"start_date" gorm:"column:start_date;type:time"`
	SignDate       *time.Time  `json:"sign_date" gorm:"column:sign_date;autoCreateTime"`
	Values         []LoanValue `json:"values,omitempty" gorm:"foreignKey:LoanId"`
}

// Transform the loan values to the loan-calculator parameters format
func (l *Loan) ToLoanCalcParameters() loancalculator.CalculationParameters {

	params := loancalculator.CalculationParameters{
		Method:         loancalculator.CalculationMethod(l.Method),
		InitialValue:   l.LoanValue,
		Rate:           l.Rate,
		RateBaseMonths: loancalculator.RateBase(l.RateBaseMonths),
		Term:           l.Term,
		BaseDate:       time.Time(l.StartDate),
	}

	return params
}

// Struct that holds the loan values data
type LoanValue struct {
	Id           *uint64   `json:"id,omitempty" gorm:"column:id;type:bigserial;autoIncrement"`
	LoanId       *string   `json:"-" gorm:"column:loan_id"`
	OwnerId      string    `json:"-" gorm:"column:owner_id"`
	Number       int       `json:"installment_number" gorm:"column:installment_number"`
	PaymentDate  LocalDate `json:"payment_date" gorm:"column:payment_date;type:time"`
	Installment  float64   `json:"installment" gorm:"column:installment"`
	Interest     float64   `json:"interest" gorm:"column:interest"`
	Amortization float64   `json:"amortization" gorm:"column:amortization"`
	Balance      float64   `json:"balance" gorm:"column:balance"`
}

// Convert the value from the loan-calculator library format
func ConvertCalculatedValue(value loancalculator.Value) LoanValue {

	return LoanValue{
		Number:       value.Number,
		PaymentDate:  LocalDate(value.PaymentDate),
		Installment:  value.Installment,
		Interest:     value.Interest,
		Amortization: value.Amortization,
		Balance:      value.Balance,
	}

}
