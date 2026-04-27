package dto

type SignInRequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponseDTO struct {
	AccessToken string `json:"access_token"`
}
