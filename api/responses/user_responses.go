package responses

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	*TokenResponse
	VerificationId string `json:"verification_id"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"success"`
}
