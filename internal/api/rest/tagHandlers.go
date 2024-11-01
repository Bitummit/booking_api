package rest

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/internal/storage/postgresql"
	"github.com/Bitummit/booking_api/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func (s *HTTPServer) ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := s.HotelService.ListTags(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, ListTagResponse{
		Tags: tags,
	})
}

func (s *HTTPServer) CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTagRequest
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

	tag := models.Tag{
		Name: req.Name,
	}
	id, err := s.HotelService.CreateTag(r.Context(), tag)
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

	s.Log.Info("New tag", slog.Int64("id", int64(id)))
	res := CreateTagResponse{
		Id: id,
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, res)
}