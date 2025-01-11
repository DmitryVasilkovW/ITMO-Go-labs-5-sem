package main

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strconv"
)

func TopCountries(w http.ResponseWriter, r *http.Request) {
	yearParam, err := parseYear(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limitParam, err := parseLimit(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filtered := filterAthletesByYear(athletes, yearParam)
	if len(filtered) == 0 {
		http.Error(w, "year not found", http.StatusNotFound)
		return
	}

	countryStats := GetCountries(filtered)
	topCountries := getTopCountries(countryStats, limitParam)

	responseJSON, err := json.Marshal(&topCountries)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, responseJSON)
}

func parseYear(r *http.Request) (int, error) {
	queries := r.URL.Query()
	yearSlc, ok := queries["year"]
	if !ok || len(yearSlc[0]) == 0 {
		return 0, http.ErrMissingFile
	}

	yearParam, err := strconv.Atoi(yearSlc[0])
	if err != nil {
		return 0, err
	}
	return yearParam, nil
}

func parseLimit(r *http.Request) (int, error) {
	queries := r.URL.Query()
	limitSlc, ok := queries["limit"]
	if ok && len(limitSlc[0]) > 0 {
		limitParam, err := strconv.Atoi(limitSlc[0])
		if err != nil {
			return 0, err
		}
		return limitParam, nil
	}
	return 3, nil
}

func filterAthletesByYear(athletes []Information, year int) []Information {
	return Filter(athletes, func(athlete Information) bool {
		return athlete.Year == year
	})
}

func getTopCountries(countries map[string]*CInfo, limit int) []*CInfo {
	vals := make([]*CInfo, 0, len(countries))
	for _, v := range countries {
		vals = append(vals, v)
	}

	sort.Slice(vals, func(i, j int) bool {
		if vals[i].Gold != vals[j].Gold {
			return vals[i].Gold > vals[j].Gold
		}
		if vals[i].Silver != vals[j].Silver {
			return vals[i].Silver > vals[j].Silver
		}
		if vals[i].Bronze != vals[j].Bronze {
			return vals[i].Bronze > vals[j].Bronze
		}
		return vals[i].Country < vals[j].Country
	})

	limit = int(math.Min(float64(limit), float64(len(vals))))
	return vals[:limit]
}

func writeJSONResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
