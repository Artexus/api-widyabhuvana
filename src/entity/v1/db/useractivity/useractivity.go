package useractivity

type status string

const (
	Completed status = "COMPLETED"
	NotYet    status = "NOT_YET"
)

type UserActivity struct {
	UserID            string `firestore:"user_id"`
	CategoryID        string `firestore:"category_id"`
	LastSubCategoryID string `firestore:"last_sub_category_id"`
	LastTaskID        string `firestore:"last_task_id"`

	Status               status `firestore:"status"`
	RemainingSubCategory int    `firestore:"remaining_sub_category"`
	RemainingTask        int    `firestore:"remaining_task"`
}
