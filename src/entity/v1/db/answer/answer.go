package answer

type Answer struct {
	TaskID string `firestore:"task_id"`
	UserID string `firestore:"user_id"`
	Answer string `firestore:"answer"`
}
