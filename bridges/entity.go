package bridges

import "time"

type Bridge struct {
	Start         time.Time `json:"start" bson:"start"`
	End           time.Time `json:"end" bson:"end"`
	HolidaysCount int       `json:"holidaysCount" bson:"holidaysCount"`
	WeekdaysCount int       `json:"weekdaysCount" bson:"weekdaysCount"`
	DaysCount     int       `json:"daysCount" bson:"daysCount"`
}

type YearBridges struct {
	Years         []string `json:"years" bson:"years"`
	Bridges       []Bridge `json:"bridges" bson:"bridges"`
	HolidaysCount int      `json:"holidaysCount" bson:"holidaysCount"`
	WeekdaysCount int      `json:"weekdaysCount" bson:"weekdaysCount"`
	DaysCount     int      `json:"daysCount" bson:"bridges"`
}
type BridgesRequest struct {
	DayOfHolidays  int              `json:"dayOfHolidays" bson:"dayOfHolidays"`
	CustomHolidays []CustomHolidays `json:"customHolidays" bson:"customHolidays"`
	City           string           `json:"city" bson:"city"`
	DaysOff        []int            `json:"daysOff" bson:"daysOff"`
	YearsScope     int              `json:"yearsScope" bson:"yearsScope"`
}

type CustomHolidays struct {
	Date string `json:"date" bson:"date"`
	Name string `json:"name" bson:"name"`
}
