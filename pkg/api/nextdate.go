package api

import (
	"fmt"
	"go_final_project/pkg/constants"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func newDay(now time.Time, start time.Time, parts []string) (string, error) {
	maxInterval := 400

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("incorrect interval: %s", parts[1])
	}

	if interval > maxInterval {
		return "", fmt.Errorf("the maximum allowed interval has been exceeded (%d):", maxInterval, interval)
	}

	date := start
	for {
		date = date.AddDate(0, 0, interval)
		if date.After(now) {
			break
		}
	}
	return date.Format(constants.DateFormat), nil
}

func newYear(now time.Time, start time.Time, parts []string) (string, error) {
	date := start
	for {
		date = date.AddDate(1, 0, 0)
		if date.After(now) {
			break
		}
	}
	return date.Format(constants.DateFormat), nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("value %s can not be empty", repeat)
	}

	start, err := time.Parse(constants.DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("failed to parse start date: %s", dstart)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) > 2 {
		return "", fmt.Errorf("repeat parameter can not be consists of more than 2 elements: %s", repeat)
	}

	switch parts[0] {
	case constants.Day:
		return newDay(now, start, parts)
	case constants.Year:
		return newYear(now, start, parts)
	case constants.Month, constants.Weekday:
		return "", fmt.Errorf("unsupported rule format %s", repeat)
	default:
		return "", fmt.Errorf("incorrect rule format: %s", repeat)
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	if dateParam == "" && repeatParam == "" {
		http.Error(w, "the date or repeat parameter is not specified", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(constants.DateFormat, nowParam)
		if err != nil {
			http.Error(w, "invalid format of the now parameter: "+nowParam, http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
