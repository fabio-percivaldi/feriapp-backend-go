package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func IsHolidays(date time.Time, daysOffMap map[int]bool, locale string, city string) bool {
	holidays := getHolidays(date.Year(), locale, city)

	return daysOffMap[int(date.Weekday())] || isCurrentDateAnHolidays(date, holidays)
}

func isCurrentDateAnHolidays(date time.Time, holidays []time.Time) bool {
	isHoliday := false

	for _, holiday := range holidays {
		if date.Equal(holiday) {
			return true
		}
	}
	return isHoliday
}

func getHolidays(year int, locale string, city string) []time.Time {
	localHolidays := readFile(locale)
	var holidays []time.Time
	var localCityHoliday Holiday
	for _, localHoliday := range localHolidays {
		if localHoliday.City == city {
			localCityHoliday = localHoliday
		}
	}
	switch locale {
	case "IT":
		holidays = []time.Time{
			time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(year, 1, 6, 0, 0, 0, 0, time.UTC),
			time.Date(year, 4, 25, 0, 0, 0, 0, time.UTC),
			time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC),
			time.Date(year, 6, 2, 0, 0, 0, 0, time.UTC),
			time.Date(year, 8, 15, 0, 0, 0, 0, time.UTC),
			time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC),
			time.Date(year, 12, 8, 0, 0, 0, 0, time.UTC),
			time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
			time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC),
		}
	default:
		return []time.Time{}
	}
	return append(holidays, localCityHoliday.Date.UTC())
}

func readFile(locale string) []Holiday {
	jsonFile, err := os.Open(fmt.Sprintf("%shelpers/%s.json", os.Getenv("LANGUAGE_PACK_FILE_PATH"), locale))
	fmt.Printf("|||||||||||| %v", jsonFile.Name())
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var holidays []Holiday

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &holidays)

	return holidays
}

type Holiday struct {
	City     string    `json:"city" bson:"city"`
	Name     string    `json:"name" bson:"name"`
	Date     time.Time `json:"date" bson:"date"`
	Region   string    `json:"region" bson:"region"`
	Province string    `json:"province" bson:"daysCount"`
}
