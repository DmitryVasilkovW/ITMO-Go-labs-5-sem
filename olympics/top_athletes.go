package main

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strconv"
)

func TopAthletes(w http.ResponseWriter, r *http.Request) {
	sportParam, limitParam, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filtered := filterAthletesBySport(athletes, sportParam)
	if len(filtered) == 0 {
		http.Error(w, "sport not found", http.StatusNotFound)
		return
	}

	filteredMap := groupAthletesByMedals(filtered)
	topAthletes := getTopAthletes(filteredMap, limitParam)

	respondWithJSON(w, topAthletes)
}

func parseQueryParams(r *http.Request) (string, int, error) {
	queries := r.URL.Query()

	sportSlc, sportOk := queries["sport"]
	if !sportOk || sportSlc[0] == "" {
		return "", 0, http.ErrMissingFile
	}

	limit := 3
	if limitSlc, limitOk := queries["limit"]; limitOk && limitSlc[0] != "" {
		var err error
		limit, err = strconv.Atoi(limitSlc[0])
		if err != nil {
			return "", 0, err
		}
	}

	return sportSlc[0], limit, nil
}

func filterAthletesBySport(athletes []Information, sport string) []Information {
	return Filter(athletes, func(athlete Information) bool {
		return athlete.Sport == sport
	})
}

func groupAthletesByMedals(filtered []Information) map[string]*AInfo {
	return GetAthlets(filtered)
}

func getTopAthletes(filteredMap map[string]*AInfo, limit int) []*AInfo {
	values := make([]*AInfo, 0, len(filteredMap))
	for _, v := range filteredMap {
		values = append(values, v)
	}

	sort.Slice(values, func(i, j int) bool {
		if values[i].Medals.Gold != values[j].Medals.Gold {
			return values[i].Medals.Gold > values[j].Medals.Gold
		}
		if values[i].Medals.Silver != values[j].Medals.Silver {
			return values[i].Medals.Silver > values[j].Medals.Silver
		}
		if values[i].Medals.Bronze != values[j].Medals.Bronze {
			return values[i].Medals.Bronze > values[j].Medals.Bronze
		}
		return values[i].Athlete < values[j].Athlete
	})

	limit = int(math.Min(float64(limit), float64(len(values))))
	return values[:limit]
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}
