package main

import (
		"fmt"
		"github.com/logpost/poc-suggestion-algorithm/utility"
		"github.com/logpost/poc-suggestion-algorithm/models"
		"github.com/logpost/poc-suggestion-algorithm/osrm"
	)

func main() {
	
	// Create OSRM client.
	osrmClient := osrm.OSRM{}
	osrmClient.CreateOSRM("http://localhost:5000/")

	// Mocking data
	jobsMock := utility.LoadJSON()

	jobMockPicked := jobsMock[1].Job
	jobsMock = jobsMock[0:len(jobsMock) - 1]
	
	// Initial data selected by user
	curentLocation := models.Location{ 
		Latitude: float64(14.7995081), 
		Longitude: float64(100.6533706),
	}

	jobPicked := models.Location{ 
		Latitude: float64(jobMockPicked.PickUpLocation.Latitude), 
		Longitude: float64(jobMockPicked.PickUpLocation.Longitude),
	}

	fmt.Println(curentLocation, jobMockPicked.PickUpLocation.Latitude, jobMockPicked.PickUpLocation.Longitude)
	
	// Initial data for running algorithm.
	res := osrmClient.GetRouteInfo(&curentLocation, &jobPicked)

	if res != nil {
		fmt.Println("responsed")
	}
	
	fmt.Printf("%+v\n",res)
	fmt.Printf("%+v\n%+v", curentLocation, jobPicked)
}