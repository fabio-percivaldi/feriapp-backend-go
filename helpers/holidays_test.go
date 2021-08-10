package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetHolidays(testCase *testing.T) {

	testCase.Run("wrong locale", func(t *testing.T) {
		actualHolidays := getHolidays(2019, "wrong_locale", "Milano")

		expectedHolidays := []time.Time{}
		require.Equal(t, expectedHolidays, actualHolidays, "Should return empty list")
	})

	testCase.Run("wrong city", func(t *testing.T) {
		year := 2019
		actualHolidays := getHolidays(year, "IT", "wrong_city")

		expectedHolidays := []time.Time{
			time.Date(year, 4, 21, 0, 0, 0, 0, time.UTC),
			time.Date(year, 4, 22, 0, 0, 0, 0, time.UTC),
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
		require.Equal(t, expectedHolidays, actualHolidays, "Should return correct list except local city holiday")
	})
}
