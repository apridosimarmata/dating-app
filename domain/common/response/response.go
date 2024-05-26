package response

import (
	"encoding/json"
	"net/http"
)

var (
	STATUS_SUCCESS = "success"
	STATUS_FAIL    = "fail"
	STATUS_ERROR   = "error"

	//TODO
	ERROR_INTERNAL_SERVER_ERROR = "Internal server error"

	ERROR_REACHED_LIMIT              = "You have reached daily acitivty limit"
	ERROR_BAD_REQUEST                = "Bad request: invalid param(s) provided"
	ERROR_UNAUTHORIZED               = "Bad credentials"
	ERROR_ANOTHER_PROCESS_IS_RUNNING = "Another process is currently running"
	ERROR_EMAIL_TAKEN                = "Email already registered"
	ERROR_USER_NOT_FOUND             = "User not found"
	ERROR_PACKAGE_NOT_FOUND          = "No such package"
	ERROR_SUBSCRIPTION_NOT_FOUND     = "Subscription not found"
)

var (
	userErrors = map[string]int{
		ERROR_REACHED_LIMIT:              429,
		ERROR_BAD_REQUEST:                400,
		ERROR_UNAUTHORIZED:               401,
		ERROR_ANOTHER_PROCESS_IS_RUNNING: 409,
		ERROR_EMAIL_TAKEN:                409,
		ERROR_USER_NOT_FOUND:             404,
		ERROR_PACKAGE_NOT_FOUND:          404,
		ERROR_SUBSCRIPTION_NOT_FOUND:     404,
	}
)

type Error struct {
	Error string `json:"error"`
}

type Response[T any] struct {
	Error   bool    `json:"error"`
	Data    *T      `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
}

func (res *Response[T]) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := 200
	if res.Error {
		statusCode = 500
		if res.Message != nil {
			if _statusCode, isUserError := userErrors[*res.Message]; isUserError {
				statusCode = _statusCode
			}
		}
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(res)
}
