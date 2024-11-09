package rest

import (
	"net/http"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/pkg/logger"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func (s *HTTPServer) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		s.Log.Error("auth: decoding request", logger.Err(err))
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

	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email: req.Email,
		FirstName: req.FirstName,
		LastName: req.LastName,
	}

	token, err := s.AuthService.Registration(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, RegistrationResponse{
		Token: token,
	})
}

func (s *HTTPServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		s.Log.Error("auth: decoding request", logger.Err(err))
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

	user := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	token, err := s.AuthService.Login(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, RegistrationResponse{
		Token: token,
	})
}