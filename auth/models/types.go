package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Date struct {
	Year  int
	Day   int
	Month int
}

//Valuer/Scanner

// syntaxis
// dd/mm/yyyy
func (d *Date) Scan(value interface{}) error {
	fmt.Println("i also clalled")
	if value == nil {
		return fmt.Errorf("Value is not valid")
	}
	split := strings.Split(fmt.Sprintf("%v", value), "/")
	day := split[0]
	month := split[1]
	year := split[2]
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return err
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return err
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return err
	}
	d.Day = dayInt
	d.Month = monthInt
	d.Year = yearInt
	return nil
}

func (d Date) Value() (driver.Value, error) {
	if d.Day == 0 || d.Year == 0 || d.Month == 0 {
		return nil, nil
	}
	return fmt.Sprintf("%v/%v/%v", d.Day, d.Month, d.Year), nil
}
