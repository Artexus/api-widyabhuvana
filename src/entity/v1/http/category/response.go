package category

type GetResponse struct {
	EncID       string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Category struct {
	EncID string `json:"id"`
	Name  string `json:"name"`
}

type SubCategory struct {
	EncID string `json:"id"`
	Name  string `json:"name"`
}

type GetUserProgressResponse struct {
	Category    `json:"category"`
	SubCategory `json:"sub_category"`

	RemainingSubCategory int `json:"remaining_sub_category"`
	TotalSubCategory     int `json:"total_sub_category"`
}
