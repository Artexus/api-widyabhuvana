package task

import "github.com/Artexus/api-widyabhuvana/src/constant"

type GetResponse struct {
	EncID string            `json:"id"`
	Type  constant.TaskType `json:"type"`
	Point int               `json:"point"`
	QnAs  []QnAPayload      `json:"qnas,omitempty"`

	Detective
	Matches Matches      `json:"matches,omitempty"`
	Levels  LevelPayload `json:"levels"`
	Learning
}

type QnaLevel struct {
	Description string `json:"description"`
	Question    string `json:"question"`
}

type Level struct {
	QnA   []QnaLevel `json:"qna"`
	Total int        `json:"total"`
}

type LevelPayload struct {
	Level1 Level `json:"level_1"`
	Level2 Level `json:"level_2"`
}

type Matches struct {
	Questions []string `json:"questions,omitempty"`
	Choices   []string `json:"choices,omitempty"`
}

type Detective struct {
	SubTasks []string `json:"sub_task"`
}

type SubmitResponse struct {
	EncID            string `json:"id,omitempty"`
	SubCategoryPoint int    `json:"sub_category_point,omitempty"`
	TotalPoint       int    `json:"total_point,omitempty"`
}

type Learning struct {
	Text     string `json:"text,omitempty"`
	VideoURL string `json:"video_url,omitempty"`
}

type QnAPayload struct {
	Question string   `json:"question,omitempty"`
	Choices  []string `json:"choices,omitempty"`
}
