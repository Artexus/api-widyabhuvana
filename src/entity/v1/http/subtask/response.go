package subtask

type QnAPayload struct {
	Question string   `json:"question"`
	Choices  []string `json:"choices"`
}

type GetResponse struct {
	EncID    string       `json:"id"`
	VideoURL string       `json:"video_url"`
	QnAs     []QnAPayload `json:"qna"`
}
