package main

import (
	"time"
	"fmt"
	"github.com/logpost/poc-suggestion-algorithm/utility"
	"github.com/logpost/poc-suggestion-algorithm/models"
	"github.com/logpost/poc-suggestion-algorithm/osrm"
)

func main() {
	
	// Create OSRM client.
	osrmClient	:=	osrm.OSRM{}
	osrmClient.CreateOSRM("http://localhost:5000/")

	// Mocking data
	jobsMock		:=	utility.LoadJSON()

	jobMockPicked	:=	jobsMock[0].Job
	// jobsMock		=	jobsMock[1:len(jobsMock) - 1]
	
	// By pass mock data to actual data
	jobs 		:=	jobsMock
	jobPicked	:=	jobMockPicked

	// Initial Priority Queue (In-Mem)
	var	Queue	[]models.Job

	// Initial variable for running algorithm
	sumCost			:=	0.0
	sumOffer		:=	0.0
	workingDays 	:=	1
	maxWorkingDays	:=	-1
	startDay	 	:=	time.Now()
	endDay			:=	time.Now()
	
	// Initial data selected by user
	curentLocation		:=	models.Location{
		Latitude:	float64(14.7995081),
		Longitude:	float64(100.6533706),
	}

	jobPickedLocation	:=	models.Location{ 
		Latitude:	float64(jobPicked.PickUpLocation.Latitude),
		Longitude:	float64(jobPicked.PickUpLocation.Longitude),
	}

	// Get cost of current location to picked job location
	routingJobPicked	:=	osrmClient.GetRouteInfo(&curentLocation, &jobPickedLocation)

	if routingJobPicked	!=	nil {
		
		preparingDistance	:=	routingJobPicked.Routes[0].Distance
		preparingCost		:=	utility.GetCostDrivingCostByDistance(preparingDistance, 0)
		sumCost				+=	preparingCost
	}

	jobPickedCost		:=	utility.GetCostDrivingCostByDistance(jobPicked.Distance, jobPicked.Weight)
	sumCost				+=	jobPickedCost
	sumOffer			+=	jobPicked.OfferPrice
	endDay				=	jobPicked.DropoffDate
	jobs[0].Job.Visited	=	true

	fmt.Printf("%+v\n%+v", curentLocation, jobPickedLocation)	
}