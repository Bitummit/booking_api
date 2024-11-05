package rest

import (
	"net/http"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/go-chi/render"
)

func (s *HTTPServer) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var req CreateHotelRequest
	// Name string 	`json:"name"`
	// 	Desc string 	`json:"desc,omitempty"`
	// 	City string 	`json:"city"`
	// 	Tags []string	`json:"tags"`
	err :=render.DecodeJSON(r.Body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse("bad request"))
		return
	}
	
	hotel := models.Hotel{
		Name: req.Name,
		Desc: req.Desc,
	}
	hotelID, err := s.HotelService.CreateHotel(r.Context(), hotel, req.City, req.Tags)
}