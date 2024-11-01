package rest

import "github.com/Bitummit/booking_api/internal/models"

type (

	Response struct{
		Status string `json:"status"`
		Error string `json:"error,omitempty"`
	}
	CreateTagRequest struct{
		Name string `json:"name"`
	}

	CreateTagResponse struct {
		Id int64 `json:"id"`
	}

	CreateCityRequest struct{
		Name string `json:"name"`
	}

	CreateCityResponse struct {
		Id int64 `json:"id"`
	}

	ListCityResponse struct {
		Cities []models.City `json:"cities"`
	}

	ListTagResponse struct {
		Tags []models.Tag `json:"tags"`
	}
)

func ErrorResponse(msg string) Response {
	return Response{
		Status: "ERROR",
		Error: msg,
	}
}
