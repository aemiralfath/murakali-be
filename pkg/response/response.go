package response

import (
	"encoding/json"
	"net/http"
)

const (
	BadRequestMessage          = "Invalid request."
	UnprocessableEntityMessage = "Request field not valid."
	InternalServerErrorMessage = "Something is wrong, pls try again later."
	NotFoundMessage            = "Route does not exist, please check again your route path."
	UnauthorizedMessage        = "Email or password not valid."
	ForbiddenMessage           = "Forbidden"

	EmailAlreadyExistMessage    = "Email already exist."
	UserAlreadyExistMessage     = "User already exist."
	UserNameAlreadyExistMessage = "Username already exist."
	PhoneNoAlreadyExistMessage  = "Phone no already exist."
	UserNotExistMessage         = "User not exist."
	OTPAlreadyExpiredMessage    = "OTP already expired."
	OTPIsNotValidMessage        = "OTP is not valid."
)

type JSONResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func returnJSONResponse(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(JSONResponse{
		Message: message,
		Data:    data,
	})
}

func SuccessResponse(w http.ResponseWriter, data interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	returnJSONResponse(w, "success", data, code)
}

func ErrorResponse(w http.ResponseWriter, message string, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	returnJSONResponse(w, message, nil, code)
}

func ErrorResponseData(w http.ResponseWriter, data interface{}, message string, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	returnJSONResponse(w, message, data, code)
}
