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
	"feriapp-backend-go/helpers"
	"net/http"
	"strconv"
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
		var reqBody bridges.BridgesRequest
		var responseBody []bridges.YearBridges

		err := json.NewDecoder(req.Body).Decode(&reqBody)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger := glogger.Get(req.Context())

		if reqBody.YearsScope == 0 {
			reqBody.YearsScope = 3
		}
		for i := 0; i < reqBody.YearsScope; i++ {
			currentYear := time.Now().UTC()
			yearBridges, err := bridgesByYear(
				currentYear.AddDate(i, 0, 0),
				4,
				reqBody.DayOfHolidays,
				reqBody.City,
				reqBody.DaysOff,
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			responseBody = append(responseBody, yearBridges)
		}

		writeResponse(logger, w, 200, responseBody)
	}
}

func bridgesByYear(date time.Time, maxHolidaysDistance int, maxAvailability int, city string, daysOff []int) (bridges.YearBridges, error) {
	var daysOffMap = make(map[int]bool)
	for i := 0; i < len(daysOff); i += 1 {
		daysOffMap[daysOff[i]] = true
	}
	startDate := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	var currentDate = startDate
	helpers.IsHolidays(currentDate, daysOffMap, "IT", city)

	var scoreMap = map[int][]bridges.Bridge{}
	var topBridges, goodBridges int

	for {
		isCurrentDateHolidays := helpers.IsHolidays(currentDate, daysOffMap, "IT", city)

		availableDays := (map[bool]int{true: maxAvailability, false: maxAvailability - 1})[isCurrentDateHolidays]
		currentBridge := bridges.Bridge{
			Start:         currentDate,
			End:           currentDate,
			HolidaysCount: (map[bool]int{true: 1, false: 0})[isCurrentDateHolidays],
			WeekdaysCount: (map[bool]int{true: 0, false: 1})[isCurrentDateHolidays],
			DaysCount:     1,
		}
		nextDate := currentDate

		for availableDays > 0 || helpers.IsHolidays(nextDate.AddDate(0, 0, 1), daysOffMap, "IT", city) {
			nextDate = nextDate.AddDate(0, 0, 1)
			isNextDateHolidays := helpers.IsHolidays(nextDate, daysOffMap, "IT", city)

			if isNextDateHolidays {
				currentBridge.HolidaysCount++
			} else {
				currentBridge.WeekdaysCount++
				availableDays -= 1
			}
			currentBridge.End = nextDate
			currentBridge.DaysCount++
		}
		currentDate = currentDate.AddDate(0, 0, 1)

		score := getBridgeScore(currentBridge)

		scoreMap[int(score)] = append(scoreMap[int(score)], currentBridge)

		if currentDate.UTC().Equal(time.Date(date.Year(), 12, 31, 0, 0, 0, 0, time.UTC)) {
			topBridges = 0
			goodBridges = 0
			for k := range scoreMap {
				if k > topBridges {
					goodBridges = topBridges
					topBridges = k
				} else {
					if k > goodBridges {
						goodBridges = k
					}
				}
			}
			break
		}
	}

	calculatedBridges := append(scoreMap[topBridges], scoreMap[goodBridges]...)
	return bridges.YearBridges{
		Years:         []string{strconv.FormatInt(int64(date.Year()), 10)},
		Bridges:       calculatedBridges,
		HolidaysCount: 6,
		WeekdaysCount: 4,
		DaysCount:     10,
	}, nil
}

func getBridgeScore(bridge bridges.Bridge) float32 {
	return (float32(bridge.DaysCount) / (float32(bridge.WeekdaysCount)) * (float32(bridge.DaysCount) / 30.0) * 100)
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
