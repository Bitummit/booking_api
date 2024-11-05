package rest

import (
	"errors"
	"net/http"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/internal/storage/postgresql"
	"github.com/Bitummit/booking_api/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func (s *HTTPServer) CreateHotelHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateHotelRequest
	// Name string 	`json:"name"`
	// 	Desc string 	`json:"desc,omitempty"`
	// 	City string 	`json:"city"`
	// 	Tags []string	`json:"tags"`
	err :=render.DecodeJSON(r.Body, &req)
	if err != nil {
		s.Log.Error("hotel: decoding request", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse("bad request"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		err = err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	hotel := models.Hotel{
		Name: req.Name,
		Desc: req.Desc,
	}
	hotelID, err := s.HotelService.CreateHotel(r.Context(), hotel, req.City, req.Tags)
	if err != nil {
		s.Log.Error("hotel:", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, postgresql.ErrorTagNotExists) {
			render.JSON(w, r, ErrorResponse("no such tag"))
			return
		}
		if errors.Is(err, postgresql.ErrorCityNotExists) {
			render.JSON(w, r, ErrorResponse("no such city"))
			return
		}
		if errors.Is(err, postgresql.ErrorExists) {
			render.JSON(w, r, ErrorResponse("hotel with this name exists!"))
			return
		}
		if errors.Is(err, postgresql.ErrorNotExists) || errors.Is(err, postgresql.ErrorInsertion){
			render.JSON(w, r, ErrorResponse("bad request"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse("internal error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	res := CreationResponse{
		Id: hotelID,
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, res)

}