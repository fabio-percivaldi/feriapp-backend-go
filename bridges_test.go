package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestBridgesRoutes(testCase *testing.T) {
	testRouter := mux.NewRouter()
	setupBridgesRouter(testRouter)

	testCase.Run("/bridges - ok", func(t *testing.T) {
		expectedResponse := "[]"

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodGet, "/bridges", nil)
		require.NoError(t, requestError, "Error creating the /bridges request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, expectedResponse, string(body), "The response body should be the expected one")
	})
}
