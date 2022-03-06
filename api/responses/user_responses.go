package responses

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"success"`
}
