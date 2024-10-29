package models

import "time"

type (
	BaseModel struct {
		CreatedAt time.Time `json:"created_at"`
		Active bool 		`json:"active"`
	}

	Tag struct {
		Id int64 		`json:"id"`
		Name string 	`json:"name"`
	}

	City struct {
		Id int64 		`json:"id"`
		Name string 	`json:"name"`
	}

	Hotel struct {
		Id int64 		`json:"id"`
		Name string 	`json:"name"`
		Desc string 	`json:"desc"`
		City_id int64 	`json:"city_id"`
		Tags []int		`json:"tags"`
	}

	RoomCategory struct {
		Id int64 		`json:"id"`
		Name string 	`json:"name"`
		Price float64 	`json:"price"`
		Capacity int64 	`json:"capacity"`
		Desc string 	`json:"desc"`
		Size int64 		`json:"size"`
		Hotel_id int64 	`json:"hotel_id"`
	}

	Room struct {
		Id int64 			`json:"id"`
		Number string 		`json:"number"`
		Category_id int64 	`json:"category_id"`
	}

	Booking struct {
		Id int64 			`json:"id"`
		EntryDate time.Time `json:"entry_date"`
		LeaveDate time.Time `json:"leave_date"`
		Price float64 		`json:"price"`
		Status string 		`json:"status"`
		GuestsCount int64 	`json:"guests_count"`
		User_id int64 		`json:"user_id"`
		Room_id int64 		`json:"room_id"`
	}
)