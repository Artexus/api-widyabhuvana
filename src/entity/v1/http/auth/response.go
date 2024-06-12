package auth

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type RegisterResponse struct {
	EncID string `json:"id"`
	Token string `json:"token"`
}
