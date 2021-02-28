package main

import (
	"time"
	"fmt"
	"container/heap"
	"github.com/logpost/poc-suggestion-algorithm/pqueue"
	"github.com/logpost/poc-suggestion-algorithm/utility"
	"github.com/logpost/poc-suggestion-algorithm/models"
	"github.com/logpost/poc-suggestion-algorithm/osrm"
)

var	osrmClient	osrm.OSRM

func main() {
	
	// Create OSRM client.
	osrmClient		=	osrm.OSRM{}
	osrmClient.CreateOSRM("http://localhost:5000/")

	// Mocking data
	jobsMock		:=	utility.LoadJSON()

	// Initial cost for each job
	for index := range jobsMock {
		jobsMock[index].Job.Cost = utility.GetDrivingCostByDistance(jobsMock[index].Job.Distance, jobsMock[index].Job.Weight)
	}
	
	jobMockPicked	:=	jobsMock[0].Job
	jobsMock		=	jobsMock[1:]

	// By pass mock data to actual data
	jobs 			:=	&jobsMock
	jobPicked		:=	jobMockPicked

	// Initial Priority Queue (In-Mem)
	var Queue	pqueue.PriorityQueue
	heap.Init(&Queue)

	// Initial variable for running algorithm
	sumCost			:=	0.0
	sumOffer		:=	0.0
	workingDays 	:=	1
	maxWorkingDays	:=	-1
	startDay	 	:=	time.Now()
	endDay			:=	time.Now()

	// Initial data selected by user
	originLocation		:=	models.CreateLocation(float64(14.7995081), float64(100.6533706))
	curentLocation		:=	originLocation

	// Starting suggestion algorithm
	item	:=	&pqueue.Item{
		Job:		&jobPicked,
		JobIndex:	0,
	}
	
	heap.Push(&Queue, item)

	for Queue.Len() > 0 {

		jobPicked	:=	heap.Pop(&Queue).(*pqueue.Item)
		(*jobs)[jobPicked.JobIndex].Job.Visited	=	true

		sumCost		+=	jobPicked.Job.Cost
		sumOffer	+=	jobPicked.Job.OfferPrice
		endDay		=	jobPicked.Job.DropoffDate

		jobPickedLocation	:=	models.CreateLocation(jobPicked.Job.PickUpLocation.Latitude, jobPicked.Job.PickUpLocation.Longitude)
		prepareRouting		:=	osrmClient.GetRouteInfo(&curentLocation, &jobPickedLocation)

		if prepareRouting	!=	nil {
			preparingDistance	:=	prepareRouting.Routes[0].Distance
			preparingCost		:=	utility.GetDrivingCostByDistance(preparingDistance, 0)
			sumCost				+=	preparingCost
		}

		minimumIndex, minimumCost	:=	getJobMinimumCost(&jobPickedLocation, &originLocation, jobs)
		fmt.Println(minimumIndex, minimumCost)

		curentLocation		=	models.CreateLocation(jobPicked.Job.DropOffLocation.Latitude, jobPicked.Job.DropOffLocation.Latitude)

	}

	fmt.Println(Queue, workingDays, maxWorkingDays, startDay, endDay)
	fmt.Println(sumCost, sumOffer, endDay)
}

func getJobMinimumCost(curentLocation *models.Location, originLocation *models.Location, jobs *[]utility.JobExpected) (int, float64) {
	
	minimumIndex	:=	0
	minimumCost		:=	9999999.999
	minimumPrepare	:=	0.0
	minimumEnd		:=	0.0

	for index	:=	range *jobs {

		if	!(*jobs)[index].Job.Visited {
			predictingPickUpLocation	:=	models.CreateLocation((*jobs)[index].Job.PickUpLocation.Latitude, (*jobs)[index].Job.PickUpLocation.Longitude)
			predictingDropOffLocation	:=	models.CreateLocation((*jobs)[index].Job.DropOffLocation.Latitude, (*jobs)[index].Job.DropOffLocation.Longitude)

			prepareRouting				:=	osrmClient.GetRouteInfo(curentLocation, &predictingPickUpLocation)
			endingRouting				:=	osrmClient.GetRouteInfo(originLocation, &predictingDropOffLocation)

			if prepareRouting != nil && endingRouting != nil {
				prepareRoutingDistance	:=	prepareRouting.Routes[0].Distance
				endingRoutingDistance	:=	endingRouting.Routes[0].Distance

				// fmt.Println(index, prepareRoutingDistance, endingRoutingDistance)

				preparingCost			:=	utility.GetDrivingCostByDistance(prepareRoutingDistance, 0)
				endingCost				:=	utility.GetDrivingCostByDistance(endingRoutingDistance, 0)
				
				sumaryPredictingCost	:=	preparingCost + (*jobs)[index].Job.Cost + endingCost

				if	minimumCost > sumaryPredictingCost {
					fmt.Println("*** MINIMUM: ", index, sumaryPredictingCost)
					minimumCost		=	sumaryPredictingCost
					minimumPrepare	=	preparingCost
					minimumEnd		=	endingCost
					minimumIndex	=	index
				}
			}
		} 
	}

	fmt.Println(minimumIndex, minimumCost, minimumPrepare, minimumEnd)
	
	return	minimumIndex, minimumCost
}