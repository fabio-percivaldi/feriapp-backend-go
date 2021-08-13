package main

import (
	"bytes"
	"encoding/json"
	"feriapp-backend-go/bridges"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestBridgesRoutes(testCase *testing.T) {
	os.Setenv("LANGUAGE_PACK_FILE_PATH", "./helpers/")

	testRouter := mux.NewRouter()
	setupBridgesRouter(testRouter)

	testCase.Run("/bridges - ok", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(bridges.BridgesRequest{
			DayOfHolidays:  2,
			CustomHolidays: []bridges.CustomHolidays{},
			City:           "Milano",
			DaysOff:        []int{0, 6},
			YearsScope:     1,
		})

		request, requestError := http.NewRequest(http.MethodPost, "/bridges", bytes.NewBuffer(requestBody))
		require.NoError(t, requestError, "Error creating the /bridges request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		var actualBridges []bridges.YearBridges

		json.Unmarshal(body, &actualBridges)
		require.NoError(t, readBodyError)
		require.Equal(t, 1, len(actualBridges), "The response body should be the expected one")
	})

	testCase.Run("/bridges - no years scope is passed", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(bridges.BridgesRequest{
			DayOfHolidays:  2,
			CustomHolidays: []bridges.CustomHolidays{},
			City:           "Milano",
			DaysOff:        []int{0, 6},
		})

		request, requestError := http.NewRequest(http.MethodPost, "/bridges", bytes.NewBuffer(requestBody))
		require.NoError(t, requestError, "Error creating the /bridges request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)

		var actualBridges []bridges.YearBridges

		json.Unmarshal(body, &actualBridges)
		require.NoError(t, readBodyError)
		require.Equal(t, 3, len(actualBridges), "The response body should be the expected one")
	})

	testCase.Run("/bridges - 2 years scope", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(bridges.BridgesRequest{
			DayOfHolidays:  2,
			CustomHolidays: []bridges.CustomHolidays{},
			City:           "Milano",
			DaysOff:        []int{0, 6},
			YearsScope:     2,
		})

		request, requestError := http.NewRequest(http.MethodPost, "/bridges", bytes.NewBuffer(requestBody))
		require.NoError(t, requestError, "Error creating the /bridges request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)

		var actualBridges []bridges.YearBridges

		json.Unmarshal(body, &actualBridges)
		require.NoError(t, readBodyError)
		require.Equal(t, 2, len(actualBridges), "The response body should be the expected one")
	})
}

func TestBridgesByYear(testCase *testing.T) {
	os.Setenv("LANGUAGE_PACK_FILE_PATH", "./helpers/")

	testCase.Run("bridgesByYear", func(t *testing.T) {
		bridgesArray := []bridges.Bridge{
			{Id: "2019-04-20-2019-04-25", IsTop: true, Start: time.Date(2019, 4, 20, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 4, 25, 0, 0, 0, 0, time.UTC), HolidaysCount: 4, WeekdaysCount: 2, DaysCount: 6},
			{Id: "2019-12-21-2019-12-26", IsTop: true, Start: time.Date(2019, 12, 21, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 12, 26, 0, 0, 0, 0, time.UTC), HolidaysCount: 4, WeekdaysCount: 2, DaysCount: 6},
			{Id: "2019-12-25-2019-12-30", IsTop: true, Start: time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 12, 30, 0, 0, 0, 0, time.UTC), HolidaysCount: 4, WeekdaysCount: 2, DaysCount: 6},
			{Id: "2019-04-25-2019-04-29", IsTop: false, Start: time.Date(2019, 4, 25, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 4, 29, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 2, DaysCount: 5},
			{Id: "2019-04-27-2019-05-01", IsTop: false, Start: time.Date(2019, 4, 27, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 2, DaysCount: 5},
			{Id: "2019-05-01-2019-05-05", IsTop: false, Start: time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 5, 5, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 2, DaysCount: 5},
			{Id: "2019-08-15-2019-08-19", IsTop: false, Start: time.Date(2019, 8, 15, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 8, 19, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 2, DaysCount: 5},
			{Id: "2019-11-01-2019-11-05", IsTop: false, Start: time.Date(2019, 11, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 11, 5, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 2, DaysCount: 5},
			{Id: "2019-12-28-2020-01-01", IsTop: false, Start: time.Date(2019, 12, 28, 0, 0, 0, 0, time.UTC), End: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 2, DaysCount: 5},
		}
		YearBridges := bridges.YearBridges{
			Years:         []string{"2019"},
			Bridges:       bridgesArray,
			HolidaysCount: 6,
			WeekdaysCount: 4,
			DaysCount:     10,
		}

		expectedResponse, _ := json.Marshal(YearBridges)

		result, err := bridgesByYear(
			time.Date(2019, 4, 19, 0, 0, 0, 0, time.UTC),
			4,
			2,
			"Milano",
			[]int{0, 6},
			false,
		)

		require.Equal(t, nil, err)

		actualResponse, _ := json.Marshal(result)
		// require.Equal(t, len(bridgesArray), len(result.Bridges), "The 2019 bridges should be %v", len(bridgesArray))
		require.Equal(t, string(expectedResponse), string(actualResponse), "The 2019 bridges should be correctly calculated")
	})

	testCase.Run("bridgesByYear - max availability = 0", func(t *testing.T) {
		bridgesArray := []bridges.Bridge{
			{Id: "2019-04-20-2019-04-22", IsTop: true, Start: time.Date(2019, 4, 20, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 4, 22, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 0, DaysCount: 3},
			{Id: "2019-11-01-2019-11-03", IsTop: true, Start: time.Date(2019, 11, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 11, 3, 0, 0, 0, 0, time.UTC), HolidaysCount: 3, WeekdaysCount: 0, DaysCount: 3},
		}
		YearBridges := bridges.YearBridges{
			Years:         []string{"2019"},
			Bridges:       bridgesArray,
			HolidaysCount: 6,
			WeekdaysCount: 4,
			DaysCount:     10,
		}

		expectedResponse, _ := json.Marshal(YearBridges)

		result, err := bridgesByYear(
			time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			4,
			0,
			"Milano",
			[]int{0, 6},
			false,
		)
		require.Equal(t, nil, err)

		actualResponse, _ := json.Marshal(result)
		require.Equal(t, string(expectedResponse), string(actualResponse), "The 2019 bridges should be 2")
	})

	testCase.Run("bridgesByYear - local city holiday - milano", func(t *testing.T) {
		expectedBridge := bridges.Bridge{
			Start:         time.Date(2020, 12, 5, 0, 0, 0, 0, time.UTC),
			End:           time.Date(2020, 12, 8, 0, 0, 0, 0, time.UTC),
			HolidaysCount: 3,
			WeekdaysCount: 0,
			DaysCount:     3,
		}

		result, err := bridgesByYear(
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			4,
			0,
			"Milano",
			[]int{0, 6},
			false,
		)
		require.Equal(t, nil, err)
		var foundBridge = false

		for _, bridge := range result.Bridges {
			if bridge.Start.Equal(expectedBridge.Start) && bridge.End.Equal(expectedBridge.End) {
				foundBridge = true
			}
		}

		require.Equal(t, true, foundBridge, "In 2020 bridges there should be san Ambrogio bridge")
	})
}
