package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/fatih/color"
	"github.com/smartech75/go-microservice-test/user-service/models"
)

// SuccessMessageResponse -> success error messageformatter
func SuccessMessageResponse(msg string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	temp := &errdata{Message: "success", Status: msg}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(temp)
}

// ErrorResponse -> error formatter
func ErrorResponse(error string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
	}
	temp := &errdata{Message: error}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(temp)
}

func SuccessItemRespond(fields interface{}, modelType string, writer http.ResponseWriter) {
	_, err := json.Marshal(fields)
	type data struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
	temp := &data{Data: fields, Message: "success"}
	if err != nil {
		ServerErrResponse(err.Error(), writer)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	switch modelType {
	case "User":
		temp.Data = fields.(models.User)
	default:
		color.Red("Invalid Model Type")
	}

	json.NewEncoder(writer).Encode(temp)
}

func ServerErrResponse(error string, writer http.ResponseWriter) {
	type servererrdata struct {
		Message string `json:"msg"`
	}
	temp := &servererrdata{Message: error}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(writer).Encode(temp)
}
