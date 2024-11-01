package rest

import "github.com/Bitummit/booking_api/internal/models"

type (

	Response struct{
		Status string `json:"status"`
		Error string `json:"error,omitempty"`
	}
	CreateTagRequest struct{
		Response Response `json:"response"`
		Name string `json:"name"`
	}

	CreateTagResponse struct {
		Response Response `json:"response"`
		Id int64 `json:"id"`
	}

	CreateCityRequest struct{
		Response Response `json:"response"`
		Name string `json:"name"`
	}

	CreateCityResponse struct {
		Response Response `json:"response"`
		Id int64 `json:"id"`
	}

	ListCityResponse struct {
		Response Response `json:"response"`
		Cities []models.City
	}

	ListTagResponse struct {
		Response Response `json:"response"`
		Tags []models.Tag
	}
)


func ErrorResponse(msg string) Response {
	return Response{
		Status: "ERROR",
		Error: msg,
	}
}
