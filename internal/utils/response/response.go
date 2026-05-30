package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Repsonse struct {
	Status string `json:"status"` //always use tags
	Error  string `json:"error"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}

func ValidationError(errs validator.ValidationErrors) Repsonse {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required": //tag that we wrote in types.go oon student struct
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field())) //err.Field() will tell us from which field the error is comming
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))

		}
	}

	return Repsonse{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}

}

func GeneralError(err error) Repsonse {
	return Repsonse{
		Status: StatusError,
		Error:  err.Error(),
	}
}
