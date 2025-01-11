package main

import (
	"fmt"
	"gitlab.com/slon/shad-go/hotelbusiness"
)

func main() {
	guests := []hotelbusiness.Guest{
		{CheckInDate: 1, CheckOutDate: 2},
	}

	loads := hotelbusiness.ComputeLoad(guests)

	for _, load := range loads {
		fmt.Println(load.StartDate, load.GuestCount)
	}
}
