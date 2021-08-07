package main

import (
	"bytes"
	"encoding/json"
	"feriapp-backend-go/bridges"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestBridgesRoutes(testCase *testing.T) {
	testRouter := mux.NewRouter()
	setupBridgesRouter(testRouter)

	testCase.Run("/bridges - ok", func(t *testing.T) {
		bridgesArray := []bridges.Bridge{}
		currentYear := time.Now().UTC().Year()
		YearBridges := []bridges.YearBridges{{
			Years:         []string{strconv.FormatInt(int64(currentYear), 10)},
			Bridges:       bridgesArray,
			HolidaysCount: 6,
			WeekdaysCount: 4,
			DaysCount:     10,
		}}

		expectedResponse, _ := json.Marshal(YearBridges)

		responseRecorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(bridges.BridgesRequest{
			DayOfHolidays:  2,
			CustomHolidays: []bridges.CustomHolidays{},
			City:           "Milan",
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
		require.NoError(t, readBodyError)
		require.Equal(t, string(expectedResponse), string(body), "The response body should be the expected one")
	})

	testCase.Run("/bridges - no years scope is passed", func(t *testing.T) {
		bridgesArray := []bridges.Bridge{}
		currentYear := time.Now().UTC().Year()
		YearBridges := []bridges.YearBridges{{
			Years:         []string{strconv.FormatInt(int64(currentYear), 10)},
			Bridges:       bridgesArray,
			HolidaysCount: 6,
			WeekdaysCount: 4,
			DaysCount:     10,
		}}

		expectedResponse, _ := json.Marshal(YearBridges)

		responseRecorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(bridges.BridgesRequest{
			DayOfHolidays:  2,
			CustomHolidays: []bridges.CustomHolidays{},
			City:           "Milan",
			DaysOff:        []int{0, 6},
		})

		request, requestError := http.NewRequest(http.MethodPost, "/bridges", bytes.NewBuffer(requestBody))
		require.NoError(t, requestError, "Error creating the /bridges request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, string(expectedResponse), string(body), "The response body should be the expected one")
	})

	testCase.Run("/bridges - 2 years scope", func(t *testing.T) {
		bridgesArray := []bridges.Bridge{}
		currentYear := time.Now().UTC()
		YearBridges := []bridges.YearBridges{{
			Years:         []string{strconv.FormatInt(int64(currentYear.Year()), 10)},
			Bridges:       bridgesArray,
			HolidaysCount: 6,
			WeekdaysCount: 4,
			DaysCount:     10,
		},
			{
				Years:         []string{strconv.FormatInt(int64(currentYear.AddDate(1, 0, 0).Year()), 10)},
				Bridges:       bridgesArray,
				HolidaysCount: 6,
				WeekdaysCount: 4,
				DaysCount:     10,
			}}

		expectedResponse, _ := json.Marshal(YearBridges)

		responseRecorder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(bridges.BridgesRequest{
			DayOfHolidays:  2,
			CustomHolidays: []bridges.CustomHolidays{},
			City:           "Milan",
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
		require.NoError(t, readBodyError)
		require.Equal(t, string(expectedResponse), string(body), "The response body should be the expected one")
	})
}
