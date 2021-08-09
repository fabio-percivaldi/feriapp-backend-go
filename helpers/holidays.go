package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func HolidaysUtils(year int, daysOffMap map[int]bool, locale string, city string) func(date time.Time) bool {
	holidays := getHolidays(year, locale, city)

	return func(date time.Time) bool {
		return daysOffMap[int(date.Weekday())] || isCurrentDateAnHolidays(date, holidays)
	}
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
	splittedDate := strings.Split(localCityHoliday.Date, "-")
	month, montErr := strconv.Atoi(splittedDate[0])
	if montErr != nil {
		fmt.Println("error parsing local city holiday month")
		return holidays
	}
	day, dayErr := strconv.Atoi(splittedDate[0])
	if dayErr != nil {
		fmt.Println("error parsing local city holiday day")
		return holidays
	}
	localCityHolidayDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	fmt.Printf("bridges %v", localCityHolidayDate)

	return append(holidays, localCityHolidayDate.UTC())
}

func readFile(locale string) []Holiday {
	jsonFile, err := os.Open(fmt.Sprintf("%shelpers/%s.json", os.Getenv("LANGUAGE_PACK_FILE_PATH"), locale))
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
	City     string `json:"city" bson:"city"`
	Name     string `json:"name" bson:"name"`
	Date     string `json:"date" bson:"date"`
	Region   string `json:"region" bson:"region"`
	Province string `json:"province" bson:"daysCount"`
}
