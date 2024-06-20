package subtask

type GetRequest struct {
	ID string `json:"-"`

	EncID string `json:"id" form:"id"`
}
