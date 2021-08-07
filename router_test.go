package main

import (
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
		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodGet, "/bridges", nil)
		require.NoError(t, requestError, "Error creating the /bridges request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")
	})
}
