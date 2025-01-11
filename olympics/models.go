package main

type Information struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type Medals struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type AInfo struct {
	Athlete      string          `json:"athlete"`
	Country      string          `json:"country"`
	Medals       Medals          `json:"medals"`
	MedalsByYear map[int]*Medals `json:"medals_by_year"`
}

type CInfo struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

func athletToInfo(a *Information) AInfo {
	var info AInfo
	info.Athlete = a.Athlete
	info.MedalsByYear = make(map[int]*Medals)
	info.Country = a.Country

	return info
}
