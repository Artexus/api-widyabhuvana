package task

type GetRequest struct {
	ID string `form:"-"`

	EncID string `form:"id"`
}

type SubmitRequest struct {
	UserID string `json:"-"`
	ID     string `json:"-"`

	EncID  string   `json:"id"`
	Answer []string `json:"answers"`
}
