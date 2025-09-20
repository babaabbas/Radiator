package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"custom_status"`
	Error  string `json:"Error"`
}

const (
	Status_Ok    = "ok"
	Status_Error = "error"
)

func GeneralError(err error) Response {
	return Response{
		Status: Status_Error,
		Error:  err.Error(),
	}
}
func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
func ValidateError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field is %s required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field is %s invalid", err.Field()))
		}
	}
	return Response{
		Status: Status_Error,
		Error:  strings.Join(errMsgs, ","),
	}
}
