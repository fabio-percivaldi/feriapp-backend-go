package bridges

import "time"

type Bridge struct {
	Start         time.Time `json:"start" bson:"start"`
	End           time.Time `json:"end" bson:"end"`
	HolidaysCount int       `json:"holidaysCount" bson:"holidaysCount"`
	WeekdaysCount int       `json:"weekdaysCount" bson:"weekdaysCount"`
	DaysCount     int       `json:"daysCount" bson:"daysCount"`
}
