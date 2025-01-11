//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return []Load{}
	}

	loads := make(map[int]int)

	loads = addArrivalAndDepartureDays(loads, guests)
	minDay, maxDay := getFullRangeOfDays(guests)

	return countTotalLoads(minDay, maxDay, loads)
}

func addArrivalAndDepartureDays(loads map[int]int, guests []Guest) map[int]int {
	for _, g := range guests {
		loads[g.CheckInDate] += 1
		loads[g.CheckOutDate] -= 1
	}

	return loads
}

func getFullRangeOfDays(guests []Guest) (int, int) {
	maxDay := -1
	minDay := guests[0].CheckInDate

	for _, g := range guests {
		if g.CheckInDate < minDay {
			minDay = g.CheckInDate
		}
		if g.CheckOutDate > maxDay {
			maxDay = g.CheckOutDate
		}
	}

	return minDay, maxDay
}

func countTotalLoads(minDay, maxDay int, loads map[int]int) []Load {
	countOfGuests := 0
	var totalLoads []Load
	for day := minDay; day <= maxDay; day++ {
		countOfGuests += getExistingCount(loads, day)

		countOfGuestsOnPreviousDay := getCountOfGuestsOnPreviousDay(totalLoads)
		totalLoads = addTotalReportForLoad(totalLoads, countOfGuests, countOfGuestsOnPreviousDay, day)
	}

	return totalLoads
}

func addTotalReportForLoad(totalLoads []Load, countOfGuests, countOfGuestsOnPreviousDay, day int) []Load {
	if isFirstIteration(totalLoads) || countOfGuests != countOfGuestsOnPreviousDay {
		totalLoads = append(totalLoads, Load{
			StartDate:  day,
			GuestCount: countOfGuests,
		})
	}

	return totalLoads
}

func isFirstIteration(totalLoads []Load) bool {
	return len(totalLoads) == 0
}

func getExistingCount(loads map[int]int, index int) int {
	if count, exists := loads[index]; exists {
		return count
	}

	return 0
}

func getCountOfGuestsOnPreviousDay(totalLoads []Load) int {
	if totalLoads != nil {
		indexOfPreviousDay := len(totalLoads) - 1
		return totalLoads[indexOfPreviousDay].GuestCount
	}

	return -1
}
