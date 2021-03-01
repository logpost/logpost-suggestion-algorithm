package main

import (
	"sync"
	"time"
	"fmt"
	"container/heap"
	"github.com/logpost/poc-suggestion-algorithm/pqueue"
	"github.com/logpost/poc-suggestion-algorithm/utility"
	"github.com/logpost/poc-suggestion-algorithm/models"
	"github.com/logpost/poc-suggestion-algorithm/osrm"
)

var	osrmClient	osrm.OSRM

// MinimumCostBuffer struct for sending to minimum pipe
type MinimumCostBuffer struct {
	minimumIndex				int
	minimumEndingCost			float64
	minimumDistanceToOrigin		float64
	minimumCost					float64
}

func getJobMinimumCost(curentLocation *models.Location, originLocation *models.Location, minimumCostPipe chan MinimumCostBuffer, wg *sync.WaitGroup, jobs *[]utility.JobExpected, startIndex int, endIndex int) {
	
	defer wg.Done()

	minimumIndex			:=	-1
	minimumCost				:=	9999999.999
	minimumPrepare			:=	0.0
	minimumEndingCost		:=	0.0
	minimumDistanceToOrigin	:=	0.0

	for index:= startIndex; index<endIndex; index++ {

		if	!(*jobs)[index].Job.Visited {
			predictingPickUpLocation	:=	models.CreateLocation((*jobs)[index].Job.PickUpLocation.Latitude,	(*jobs)[index].Job.PickUpLocation.Longitude)
			predictingDropOffLocation	:=	models.CreateLocation((*jobs)[index].Job.DropOffLocation.Latitude,	(*jobs)[index].Job.DropOffLocation.Longitude)

			prepareRouting				:=	osrmClient.GetRouteInfo(curentLocation,	&predictingPickUpLocation)
			endingRouting				:=	osrmClient.GetRouteInfo(&predictingDropOffLocation,	originLocation)

			if prepareRouting != nil && endingRouting != nil {
				prepareRoutingDistance	:=	prepareRouting.Routes[0].Distance
				endingRoutingDistance	:=	endingRouting.Routes[0].Distance

				preparingCost			:=	utility.GetDrivingCostByDistance(prepareRoutingDistance, 0)
				endingCost				:=	utility.GetDrivingCostByDistance(endingRoutingDistance, 0)
				
				sumaryPredictingCost	:=	preparingCost + (*jobs)[index].Job.Cost + endingCost

				if	minimumCost > sumaryPredictingCost {
					minimumCost				=	sumaryPredictingCost
					minimumPrepare			=	preparingCost
					minimumDistanceToOrigin	=	endingRoutingDistance
					minimumEndingCost		=	endingCost
					minimumIndex			=	index
				}
			}
		} 
	}

	fmt.Printf("\nMINIMUM PREDICT:\nINDEX:\t\t%d\t\tCOST:\t\t%f\nCOST_PREPARE:\t%f\tCOST_ENDING:\t%f", minimumIndex, minimumCost, minimumPrepare, minimumEndingCost)
	 
	buffer	:=	MinimumCostBuffer{
		minimumIndex, minimumEndingCost, minimumDistanceToOrigin, minimumCost,
	}

	minimumCostPipe	<- buffer

}

func getActualJobMinimumCost(minimumCostPipe chan MinimumCostBuffer, actualJobMinimumCostPipe chan MinimumCostBuffer, wg *sync.WaitGroup) {

	defer wg.Done()

	var minimumCostOne,	minimumCostTwo	MinimumCostBuffer
	m, mm, mmm, mmmm := <-minimumCostPipe, <-minimumCostPipe, <-minimumCostPipe, <-minimumCostPipe

	if m.minimumCost > mm.minimumCost {
		minimumCostOne	=	m
	} else {
		minimumCostOne	=	mm
	}

	if mmm.minimumCost > mmmm.minimumCost {
		minimumCostTwo	=	mmm
	} else {
		minimumCostTwo	=	mmmm
	}
	
	if minimumCostOne.minimumCost < minimumCostTwo.minimumCost {
		actualJobMinimumCostPipe <- minimumCostOne
	} else {
		actualJobMinimumCostPipe <- minimumCostTwo
	}
}

func getJobMinimumCostParallel(jobPickedLocation *models.Location, originLocation *models.Location, jobs *[]utility.JobExpected) (int, float64, float64, float64) {

	minimumCostPipe				:=	make(chan MinimumCostBuffer)
	actualJobMinimumCostPipe	:=	make(chan MinimumCostBuffer)

	var wg sync.WaitGroup

	wg.Add(1)
	go getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &wg, jobs, 0, len(*jobs)/4)

	wg.Add(1)
	go getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &wg, jobs, len(*jobs)/4, len(*jobs)/4 * 2)

	wg.Add(1)
	go getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &wg, jobs, len(*jobs)/4 * 2, len(*jobs)/4 * 3)

	wg.Add(1)
	go getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &wg, jobs, len(*jobs)/4 * 3, len(*jobs))

	wg.Add(1)
	go getActualJobMinimumCost(minimumCostPipe, actualJobMinimumCostPipe, &wg)

	actualJobMinimumCost		:=	<- actualJobMinimumCostPipe
	minimumIndex				:=	actualJobMinimumCost.minimumIndex
	minimumEndingCost			:=	actualJobMinimumCost.minimumEndingCost
	minimumDistanceToOrigin		:=	actualJobMinimumCost.minimumDistanceToOrigin
	minimumCost					:=	actualJobMinimumCost.minimumCost
	
	wg.Wait()

	return minimumIndex, minimumEndingCost, minimumDistanceToOrigin, minimumCost
	
}

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
	
	jobMockPicked	:=	jobsMock[30].Job
	jobsMock		=	jobsMock[1:]

	fmt.Println(jobMockPicked)

	// By pass mock data to actual data
	jobs	 		:=	&jobsMock
	jobPicked		:=	jobMockPicked

	// Initial Priority Queue (In-Mem)
	var minimumIndex			int
	var minimumEndingCost		float64
	var minimumDistanceToOrigin	float64
	var Queue					pqueue.PriorityQueue
	
	heap.Init(&Queue)

	// Initial variable for running algorithm
	sumCost			:=	0.0
	sumOffer		:=	0.0
	currentHop		:=	0
	maxHop			:=	20
	workingDays 	:=	1
	maxWorkingDays	:=	-1
	startDay	 	:=	time.Now()
	endDay			:=	time.Now()

	// Initial data selected by user
	originLocation	:=	models.CreateLocation(float64(14.7995081), float64(100.6533706))
	curentLocation	:=	originLocation

	// Starting suggestion algorithm
	majorJob	:=	&pqueue.Item{
		Job:		&jobPicked,
		JobIndex:	0,
	}
	
	heap.Push(&Queue, majorJob)

	for Queue.Len() > 0 {
		
		currentHop++
		
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

		if currentHop	<=	maxHop {
			
			minimumIndex, minimumEndingCost, minimumDistanceToOrigin, _	=	getJobMinimumCostParallel(&jobPickedLocation, &originLocation, jobs)

			if minimumIndex	!=	-1 {
				
				job	:=	&pqueue.Item{
					Job:		&(*jobs)[minimumIndex].Job,
					JobIndex:	minimumIndex,
				}

				heap.Push(&Queue, job)
				curentLocation	=	models.CreateLocation((*jobs)[minimumIndex].Job.DropOffLocation.Latitude, (*jobs)[minimumIndex].Job.DropOffLocation.Latitude)

			}
		}

		if currentHop > maxHop	||	Queue.Len() == 0 {
			
			if currentHop == 1 {
				predictingDropOffLocation	:=	models.CreateLocation(jobPicked.Job.PickUpLocation.Latitude, jobPicked.Job.PickUpLocation.Longitude)
				endingRouting				:=	osrmClient.GetRouteInfo(&predictingDropOffLocation, &originLocation)

				if endingRouting	!=	nil {
					endingRoutingDistance	:=	endingRouting.Routes[0].Distance
					minimumDistanceToOrigin	=	endingRoutingDistance
					endingCost				:=	utility.GetDrivingCostByDistance(endingRoutingDistance, 0)
					sumCost					+=	endingCost
				}

			} else {
				sumCost	+=	minimumEndingCost
			}

			break

		}

		fmt.Println("\nCURRENT_HOP: ", currentHop)

	}
	
	fmt.Printf("\n## SUMARY ##\n")
	fmt.Printf("SUM_OFFER:\t\t%f\n",		sumOffer)
	fmt.Printf("SUM_COST:\t\t%f\n",			sumCost)
	fmt.Printf("SUM_PROFIT:\t\t%f\n",		sumOffer - sumCost)
	fmt.Printf("START_DATE:\t\t%s\n",		startDay.String())
	fmt.Printf("END_DATE:\t\t%s\n",			endDay.String())
	fmt.Printf("DISTANCE_TO_ORIGIN:\t%f\n",	minimumDistanceToOrigin)

	fmt.Println("DEBUG: ", Queue, workingDays, maxWorkingDays, startDay, endDay)
}
