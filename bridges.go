/*
 * Copyright 2019 Mia srl
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"errors"
	"feriapp-backend-go/bridges"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mia-platform/glogger"
	"github.com/sirupsen/logrus"
)

var (
	errGeneric = errors.New("internal Server Error")
	// errBadRequest = errors.New("bad Request")
)

func setupBridgesRouter(router *mux.Router) {
	// Setup your routes here.
	router.HandleFunc("/bridges", createBridges()).Methods(http.MethodPost)
}

func createBridges() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logger := glogger.Get(req.Context())

		writeResponse(logger, w, 200, [2]bridges.Bridge{
			{Start: time.Date(2019, 12, 21, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 12, 26, 0, 0, 0, 0, time.UTC), HolidaysCount: 4, WeekdaysCount: 2, DaysCount: 6},
			{Start: time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC), End: time.Date(2019, 12, 29, 0, 0, 0, 0, time.UTC), HolidaysCount: 4, WeekdaysCount: 1, DaysCount: 5},
		})
	}
}

func writeResponse(logger *logrus.Entry, w http.ResponseWriter, statusCode int, response interface{}) {
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.WithError(err).Error("failed response unmarshalling")
		http.Error(w, errGeneric.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(responseBody)
}
