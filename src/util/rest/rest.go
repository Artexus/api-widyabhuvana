package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
}

func ResponseOutput(ctx *gin.Context, status int, detail interface{}) {
	response := Response{
		Status: http.StatusText(status),
	}
	if detail != nil {
		response.Detail = detail
		if err, ok := detail.(error); ok {
			response.Detail = map[string]string{
				"error": err.Error(),
			}
		}
	}

	ctx.JSON(status, response)
}

func ResponseData(ctx *gin.Context, status int, entity interface{}) {
	response := Response{
		Status: http.StatusText(status),
		Detail: entity,
	}
	ctx.JSON(status, response)
}
