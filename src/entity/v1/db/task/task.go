package task

import (
	"encoding/json"

	"github.com/Artexus/api-widyabhuvana/src/constant"
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
)

type Task struct {
	ID            string            `firestore:"id"`
	Type          constant.TaskType `firestore:"type"`
	Point         int               `firestore:"point"`
	CategoryID    string            `firestore:"category_id"`
	SubCategoryID string            `firestore:"sub_category_id"`
	SubTasks      []string          `firestore:"sub_tasks"`

	Level interface{} `firestore:"levels"`
	Learning
	QnA
}

type QnaLevel struct {
	Answer      []string `json:"answer"`
	Description string   `json:"description"`
	Question    string   `json:"question"`
}

type Level struct {
	QnA   []QnaLevel `json:"qna"`
	Total int        `json:"total"`
}

type LevelPayload struct {
	Level1 Level `json:"level_1"`
	Level2 Level `json:"level_2"`
}

type Learning struct {
	Text     string `firestore:"text"`
	VideoURL string `firestore:"video_url"`
}

type QnA struct {
	QnAs []interface{} `firestore:"qna"`
}

type QnAPayload struct {
	Question string   `firestore:"question"`
	Choices  []string `firestore:"choices"`
	Answer   string   `firestore:"answer"`
}

func (t Task) EncID() string {
	return aes.EncryptID(t.ID)
}

func (mc QnA) Payload() []QnAPayload {
	p := []QnAPayload{}
	r, _ := json.Marshal(mc.QnAs)

	json.Unmarshal(r, &p)
	return p
}

func (t Task) LevelPayload() LevelPayload {
	p := LevelPayload{}
	r, _ := json.Marshal(t.Level)

	json.Unmarshal(r, &p)
	return p
}
