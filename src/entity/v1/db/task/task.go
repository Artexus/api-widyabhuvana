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

	Learning
	MultipleChoice
}

type Learning struct {
	Text     string `firestore:"text"`
	VideoURL string `firestore:"video_url"`
}

type MultipleChoice struct {
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

func (mc MultipleChoice) Payload() []QnAPayload {
	p := []QnAPayload{}
	r, _ := json.Marshal(mc.QnAs)

	json.Unmarshal(r, &p)
	return p
}
