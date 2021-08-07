package main

import (
	"encoding/json"
	"feriapp-backend-go/bridges"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestBridgesRoutes(testCase *testing.T) {
	testRouter := mux.NewRouter()
	setupBridgesRouter(testRouter)

	testCase.Run("/bridges - ok", func(t *testing.T) {
		bridgesArray := [2]bridges.Bridge{
			{Start: time.Date(2019, 12, 21, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 12, 26, 0, 0, 0, 0, time.UTC), HolidaysCount: 4, WeekdaysCount: 2, DaysCount: 6},
			{Start: time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 12, 29, 0, 0, 0, 0, time.UTC), HolidaysCount: 4, WeekdaysCount: 1, DaysCount: 5},
		}
		fmt.Println(bridgesArray)
		expectedResponse, _ := json.Marshal(bridgesArray)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodPost, "/bridges", nil)
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
