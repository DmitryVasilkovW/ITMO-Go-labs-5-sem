package main

func Filter(list []Information, predicate func(Information) bool) (ret []Information) {
	for _, athlete := range list {
		if predicate(athlete) {
			ret = append(ret, athlete)
		}
	}
	return
}

func GetAthlets(slc []Information) map[string]*AInfo {
	all := make(map[string]*AInfo)

	for _, athlete := range slc {
		_, ok := all[athlete.Athlete]
		if !ok {
			medalsByYear := make(map[int]*Medals)
			var info AInfo
			info.Athlete = athlete.Athlete
			info.MedalsByYear = medalsByYear
			info.Country = athlete.Country

			all[athlete.Athlete] = &info
		}

		info := all[athlete.Athlete]

		_, ok2 := info.MedalsByYear[athlete.Year]
		if !ok2 {
			info.MedalsByYear[athlete.Year] = &Medals{0, 0, 0, 0}
		}
		info.MedalsByYear[athlete.Year].Gold += athlete.Gold
		info.MedalsByYear[athlete.Year].Silver += athlete.Silver
		info.MedalsByYear[athlete.Year].Bronze += athlete.Bronze
		info.MedalsByYear[athlete.Year].Total += athlete.Gold + athlete.Silver + athlete.Bronze

		info.Medals.Gold += athlete.Gold
		info.Medals.Silver += athlete.Silver
		info.Medals.Bronze += athlete.Bronze
		info.Medals.Total += athlete.Gold + athlete.Silver + athlete.Bronze
	}

	return all
}

func GetCountries(slc []Information) map[string]*CInfo {
	all := make(map[string]*CInfo)

	for _, athlete := range slc {
		_, ok := all[athlete.Country]
		if !ok {
			s := CInfo{athlete.Country, 0, 0, 0, 0}
			all[athlete.Country] = &s
		}

		all[athlete.Country].Gold += athlete.Gold
		all[athlete.Country].Silver += athlete.Silver
		all[athlete.Country].Bronze += athlete.Bronze
		all[athlete.Country].Total += athlete.Gold + athlete.Silver + athlete.Bronze
	}

	return all
}
