package endpoints

type ErrorResponse struct {
	Message string `json:"message"`
}

func ErrorResponseFromError(err error) ErrorResponse {
	return ErrorResponse{Message: err.Error()}
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Message: message}
}
