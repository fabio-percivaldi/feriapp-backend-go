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
	"fmt"
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
	router.HandleFunc("/bridges", createBridges).Methods(http.MethodPost)
}

func createBridges(w http.ResponseWriter, req *http.Request) {
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
			true,
		)
		filteredBridges := []bridges.Bridge{}
		for _, bridge := range yearBridges.Bridges {
			if bridge.Start.After(time.Now().UTC().AddDate(0, 0, reqBody.DayOfHolidays)) {
				filteredBridges = append(filteredBridges, bridge)
			}
		}
		yearBridges.Bridges = filteredBridges
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		responseBody = append(responseBody, yearBridges)
	}

	writeResponse(logger, w, 200, responseBody)
}

func bridgesByYear(date time.Time, maxHolidaysDistance int, maxAvailability int, city string, daysOff []int, skipPastBridges bool) (bridges.YearBridges, error) {
	var daysOffMap = make(map[int]bool)
	for i := 0; i < len(daysOff); i += 1 {
		daysOffMap[daysOff[i]] = true
	}
	startDate := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	var currentDate = startDate

	var isHolidays = helpers.HolidaysUtils(currentDate.Year(), daysOffMap, "IT", city)

	var scoreMap = map[int][]bridges.Bridge{}
	var topBridges, goodBridges int
	var calculatedBridges []bridges.Bridge

	for !currentDate.UTC().After(time.Date(date.Year(), 12, 31, 0, 0, 0, 0, time.UTC)) {

		isCurrentDateHolidays := isHolidays(currentDate)
		availableDays := (map[bool]int{true: maxAvailability, false: maxAvailability - 1})[isCurrentDateHolidays]
		// if no more days off are left and today is not holiday the bridge is closed
		if maxAvailability == 0 && !isCurrentDateHolidays {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}
		// if skipPastBridges is true only bridges that happens after today - maxAvailability day will be returned
		if currentDate.Before(time.Now().AddDate(0, 0, -(maxAvailability+1))) && skipPastBridges {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}
		currentBridge := bridges.Bridge{
			Start:         currentDate,
			End:           currentDate,
			HolidaysCount: (map[bool]int{true: 1, false: 0})[isCurrentDateHolidays],
			WeekdaysCount: (map[bool]int{true: 0, false: 1})[isCurrentDateHolidays],
			DaysCount:     1,
		}

		nextDate := currentDate
		nextDate = nextDate.AddDate(0, 0, 1)
		// a bridge should always start with an holiday
		if currentBridge.DaysCount == 1 && !isCurrentDateHolidays {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		for availableDays > 0 || isHolidays(nextDate) {
			isNextDateHolidays := isHolidays(nextDate)

			if isNextDateHolidays {
				currentBridge.HolidaysCount++
			} else {
				currentBridge.WeekdaysCount++
				availableDays -= 1
			}

			currentBridge.End = nextDate
			currentBridge.DaysCount++
			nextDate = nextDate.AddDate(0, 0, 1)
			if nextDate.Year() != currentDate.Year() {
				isHolidays = helpers.HolidaysUtils(nextDate.Year(), daysOffMap, "IT", city)
			}
		}
		for isHolidays(currentDate) {
			currentDate = currentDate.AddDate(0, 0, 1)
		}

		score := getBridgeScore(currentBridge)
		// the bridge is inserted only if it is longer than daysOff (es: exlude weekend bridges)
		// and if it is not in the past for more than maxAvailability days
		currentBridge.Id = fmt.Sprintf("%s-%s", currentBridge.Start.Format("2006-01-02"), currentBridge.End.Format("2006-01-02"))

		if currentBridge.DaysCount > len(daysOff) {
			scoreMap[int(score)] = append(scoreMap[int(score)], currentBridge)
		}
	}
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

	for index := range scoreMap[topBridges] {
		scoreMap[topBridges][index].IsTop = true
	}

	calculatedBridges = append(scoreMap[topBridges], scoreMap[goodBridges]...)

	return bridges.YearBridges{
		Years:         []string{strconv.FormatInt(int64(date.Year()), 10)},
		Bridges:       calculatedBridges,
		HolidaysCount: 6,
		WeekdaysCount: 4,
		DaysCount:     10,
	}, nil
}

func getBridgeScore(bridge bridges.Bridge) float32 {
	if bridge.WeekdaysCount == 0 {
		return float32(bridge.DaysCount)
	}
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
