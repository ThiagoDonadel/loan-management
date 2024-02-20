package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Type used customize the time format
type LocalDate time.Time

func (l LocalDate) MarshalJSON() ([]byte, error) {
	formatedDate := time.Time(l).Format(time.DateOnly)
	return json.Marshal(formatedDate)
	//return []byte(formatedDate), nil
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
