package task

import "github.com/Artexus/api-widyabhuvana/src/constant"

type GetResponse struct {
	EncID string            `json:"id"`
	Type  constant.TaskType `json:"type"`
	Point int               `json:"point"`
	QnAs  []QnAPayload      `json:"qnas,omitempty"`

	Learning
}

type Learning struct {
	Text     string `json:"text,omitempty"`
	VideoURL string `json:"video_url,omitempty"`
}

type QnAPayload struct {
	Question string   `json:"question,omitempty"`
	Choices  []string `json:"choices,omitempty"`
}
