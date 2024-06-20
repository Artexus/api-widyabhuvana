package user

type GetResponse struct {
	EncID      string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	TotalPoint int64  `json:"total_point"`
	PhotoURL   string `json:"photo_url"`
}
