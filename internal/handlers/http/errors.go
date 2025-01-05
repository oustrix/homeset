package http

import (
	"net/http"
	"strings"

	"github.com/oustrix/homeset/internal/pkg/homeset/http/api"
)

func ErrorHandler(w http.ResponseWriter, message string, statusCode int) {
	switch statusCode {
	case http.StatusNotFound:
		responseError(w, statusCode, "path not found")
	case http.StatusBadRequest:
		handleBadRequestError(w, message)
	default:
		responseError(w, statusCode, message)
	}
}

func handleBadRequestError(w http.ResponseWriter, message string) {
	if !strings.Contains(message, "doesn't match schema") {
		responseError(w, http.StatusBadRequest, message)
		return
	}

	parts := strings.Split(message, ":")
	if len(parts) != 4 {
		responseError(w, http.StatusBadRequest, message)
		return
	}

	errorMessage := strings.Join(parts[2:], ":")

	errorMessage = strings.TrimSpace(errorMessage)
	errorMessage = strings.Replace(errorMessage, "/", "", -1)

	responseError(w, http.StatusBadRequest, errorMessage)
}

func responseError(w http.ResponseWriter, code int, message string) {
	apiErr := api.Error{
		StatusCode: code,
		Error:      message,
	}

	response(w, code, apiErr)
}
