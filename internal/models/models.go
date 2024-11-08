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
		CityId int64 	`json:"city_id"`
	}

	RoomCategory struct {
		Id int64 		`json:"id"`
		Name string 	`json:"name"`
		Price float64 	`json:"price"`
		Capacity int64 	`json:"capacity"`
		Desc string 	`json:"desc"`
		Size int64 		`json:"size"`
		HotelId int64 	`json:"hotel_id"`
	}

	Room struct {
		Id int64 			`json:"id"`
		Number string 		`json:"number"`
		CategoryId int64 	`json:"category_id"`
	}

	Booking struct {
		Id int64 			`json:"id"`
		EntryDate time.Time `json:"entry_date"`
		LeaveDate time.Time `json:"leave_date"`
		Price float64 		`json:"price"`
		Status string 		`json:"status"`
		GuestsCount int64 	`json:"guests_count"`
		UserId int64 		`json:"user_id"`
		RoomId int64 		`json:"room_id"`
	}

	HotelTag struct {
		Id int64 			`json:"id"`
		CityId string 		`json:"city_id"`
		Tagid int64 		`json:"tag_id"`
	}

	User struct {
		Id int64
		Username string
		FirstName string
		LastName string
		Birthday time.Time
		Email string
		Password string
		Role string
	}
)