package models

import (
	"time"
)

// Job Struct for create job instance
type Job struct { 
	OfferPrice 			int 		`json:"offer_price"`
	Weight				int 		`json:"weight"`
	Duration 			int 		`json:"duration"`
	Distance 			float64 	`json:"distance"`
	ProductType			string 		`json:"product_type"`
	Permission			string 		`json:"permission"`
	PickupDate			time.Time 	`json:"pickup_date"`
	DropoffDate			time.Time 	`json:"dropoff_date"`
	PickUpLocation		Location	`json:"pickup_location"`
	DropOffLocation		Location	`json:"dropoff_location"`
}

// Location Struct for mapping location information
type Location struct {
	Latitude	float64			`json:"latitude"`
	Longitude	float64			`json:"longitude"`
	Address		string			`json:"address"`
	Province	string			`json:"province"`
	District	string			`json:"district"`
	Zipcode		string			`json:"zipcode"`
}