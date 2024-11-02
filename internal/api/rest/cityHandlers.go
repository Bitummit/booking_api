package rest

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/internal/storage/postgresql"
	"github.com/Bitummit/booking_api/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func (s *HTTPServer) ListCityHandler(w http.ResponseWriter, r *http.Request) {
	cities, err := s.HotelService.ListCities(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, ListCityResponse{
		Cities: cities,
	})
}

func (s *HTTPServer) CreateCityHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateCityRequest
	err := render.DecodeJSON(r.Body, &req)
	r.Body.Close()
	if err != nil {
		s.Log.Error("decoding request: %v", logger.Err(err))
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

	city := models.City{
		Name: req.Name,
	}
	id, err := s.HotelService.CreateCity(r.Context(), city)
	if err != nil {
		s.Log.Error("%v", logger.Err(err))
		if errors.Is(err, postgresql.ErrorInsertion){
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse("insertion error"))
			return
		} else if errors.Is(err, postgresql.ErrorExists){
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse("tag aready exists"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	s.Log.Info("New city", slog.Int64("id", int64(id)))
	res := CreationResponse{
		Id: id,
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, res)
}

func (s *HTTPServer) DeleteCityHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse("id is not int"))
		return
	}

	err = s.HotelService.DeleteCity(r.Context(), int64(id))
	if err != nil {
		s.Log.Error("deleting city ", logger.Err(err))
		if errors.Is(err, postgresql.ErrorNotExists) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrorResponse("city not exists"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse("internal error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{Status: "OK"})
}