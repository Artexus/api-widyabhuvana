package subcategory

type GetRequest struct {
	UserID     string `json:"-"`
	CategoryID string `json:"-"`

	EncCategoryID string `json:"category_id" form:"category_id"`
}

type GetResponse struct {
	EncID     string   `json:"id"`
	Name      string   `json:"name"`
	MaxPoint  int      `json:"max_point"`
	Tasks     []string `json:"tasks"`
	TotalTask int      `json:"total_task"`
}
