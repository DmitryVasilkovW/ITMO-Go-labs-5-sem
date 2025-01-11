package main

import (
	"encoding/json"
	"net/http"
)

func Athletes(w http.ResponseWriter, r *http.Request) {
	name, err := getNameFromQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filtered := Filter(athletes, func(athlete Information) bool {
		return athlete.Athlete == name
	})

	if len(filtered) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := aggregateMedals(filtered)

	err = sendResponse(w, resp)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func getNameFromQuery(r *http.Request) (string, error) {
	queries := r.URL.Query()
	if nameSlc, ok := queries["name"]; !ok || len(nameSlc[0]) <= 0 {
		return "", http.ErrNoLocation
	}
	nameSlc := queries["name"]
	return nameSlc[0], nil
}

func aggregateMedals(filtered []Information) AInfo {
	resp := athletToInfo(&filtered[0])

	for _, athlete := range filtered {
		if _, ok := resp.MedalsByYear[athlete.Year]; !ok {
			resp.MedalsByYear[athlete.Year] = &Medals{0, 0, 0, 0}
		}
		resp.MedalsByYear[athlete.Year].Gold += athlete.Gold
		resp.MedalsByYear[athlete.Year].Silver += athlete.Silver
		resp.MedalsByYear[athlete.Year].Bronze += athlete.Bronze
		resp.MedalsByYear[athlete.Year].Total += athlete.Gold + athlete.Silver + athlete.Bronze

		resp.Medals.Gold += athlete.Gold
		resp.Medals.Silver += athlete.Silver
		resp.Medals.Bronze += athlete.Bronze
		resp.Medals.Total += athlete.Gold + athlete.Silver + athlete.Bronze
	}

	return resp
}

func sendResponse(w http.ResponseWriter, resp AInfo) error {
	b, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	return err
}
