package loan

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	loancalculator "github.com/ThiagoDonadel/loan-calculator"
)

var (
	ErrInvalidDateFormat = errors.New("invalid date format")
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

type Loan struct {
	Id             *string     `json:"id" gorm:"column:id;primarykey;default:gen_random_uuid()"`
	Method         int         `json:"method" gorm:"column:method"`
	LoanValue      float64     `json:"value" gorm:"column:loan_value"`
	Rate           float64     `json:"rate" gorm:"column:rate"`
	RateBaseMonths int         `json:"rate_base_months" gorm:"column:rate_base_months"`
	Term           int         `json:"term" gorm:"column:term"`
	StartDate      LocalDate   `json:"start_date" gorm:"column:start_date;type:time"`
	SignDate       *time.Time  `json:"sign_date" gorm:"column:sign_date;autoCreateTime"`
	Values         []LoanValue `json:"values,omitempty" gorm:"foreignKey:LoanId"`
}

func (l *Loan) toLoanCalcParameters() loancalculator.CalculationParameters {

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

type LoanValue struct {
	Id           *uint64   `json:"id,omitempty" gorm:"column:id;type:bigserial;autoIncrement"`
	LoanId       *string   `json:"-" gorm:"column:loan_id"`
	Number       int       `json:"installment_number" gorm:"column:installment_number"`
	PaymentDate  LocalDate `json:"payment_date" gorm:"column:payment_date;type:time"`
	Installment  float64   `json:"installment" gorm:"column:installment"`
	Interest     float64   `json:"interest" gorm:"column:interest"`
	Amortization float64   `json:"amortization" gorm:"column:amortization"`
	Balance      float64   `json:"balance" gorm:"column:balance"`
	Paid         *bool     `json:"paid,omitempty" gorm:"column:paid;default:false"`
}

func convertCalculatedValue(value loancalculator.Value) LoanValue {

	return LoanValue{
		Number:       value.Number,
		PaymentDate:  LocalDate(value.PaymentDate),
		Installment:  value.Installment,
		Interest:     value.Interest,
		Amortization: value.Amortization,
		Balance:      value.Balance,
	}

}
