package subtask

import (
	"encoding/json"

	"github.com/Artexus/api-widyabhuvana/src/util/aes"
)

type SubTask struct {
	ID       string      `firestore:"id"`
	QnAs     interface{} `firestore:"qna"`
	VideoURL string      `firestore:"video_url"`
}

type QnA struct {
	Answer   []string `json:"answer"`
	Choices  []string `json:"choices"`
	Question string   `json:"question"`
}

func (mc SubTask) Payload() []QnA {
	p := []QnA{}
	r, _ := json.Marshal(mc.QnAs)

	json.Unmarshal(r, &p)
	return p
}

func (st SubTask) EncID() string {
	return aes.EncryptID(st.ID)
}
