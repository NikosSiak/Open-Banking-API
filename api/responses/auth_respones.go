package responses

type UnauthorizedError struct {
	Message string `json:"message" example:"Missing or Invalid token"`
}
