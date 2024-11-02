package rest

import "github.com/Bitummit/booking_api/internal/models"

type (
	Response struct{
		Status string `json:"status"`
		Error string `json:"error,omitempty"`
	}
	CreationResponse struct {
		Id int64 `json:"id"`
	}
	CreateTagRequest struct{
		Name string `json:"name"`
	}
	CreateCityRequest struct{
		Name string `json:"name"`
	}
	ListCityResponse struct {
		Cities []models.City `json:"cities"`
	}
	ListTagResponse struct {
		Tags []models.Tag `json:"tags"`
	}
	CreateHotelRequest struct{
		Name string 	`json:"name"`
		Desc string 	`json:"desc,omitempty"`
		City string 	`json:"city"`
		Tags []string	`json:"tags"`
	}
)

func ErrorResponse(msg string) Response {
	return Response{
		Status: "ERROR",
		Error: msg,
	}
}
