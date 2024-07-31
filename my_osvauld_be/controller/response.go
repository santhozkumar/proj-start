package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ResponseStatus string

const (
	ResponseSuccess ResponseStatus = "success"
	ResponseFail    ResponseStatus = "fail"
	ResponseError   ResponseStatus = "error"
)

type ResponseCode int

const (
	FailResponseCode ResponseCode = 800
)

type APIResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
	Code    ResponseCode   `json:"code,omitempty"`
}

func SendResponse(c *gin.Context, httpCode int, status ResponseStatus,
	data interface{}, message string, code ResponseCode, err error) {

	// jsend, json reponse format https://github.com/omniti-labs/jsend
	var response APIResponse

	switch status {
	case ResponseSuccess:
		response = APIResponse{
			Status:  status,
			Message: message, // optional
			Data:    data,    // Required
		}
		c.JSON(httpCode, response)
	case ResponseError:
		response = APIResponse{
			Status:  status,
			Message: err.Error(),
			Code:    code, // optional
		}
		c.JSON(httpCode, response)
	case ResponseFail:
		// pending, not needed for most cases
		response = APIResponse{
			Status:  status,
			Message: err.Error(),
			Code:    code, // optional
		}
		c.JSON(httpCode, response)
	default:
		fmt.Println("Unknown outcome.")
	}
}
