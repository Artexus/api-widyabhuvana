package subcategory

import "github.com/Artexus/api-widyabhuvana/src/util/aes"

type SubCategory struct {
	ID         string   `firestore:"id"`
	CategoryID string   `firestore:"category_id"`
	Name       string   `firestore:"name"`
	Tasks      []string `firestore:"tasks"`
	MaxPoint   int      `firestore:"max_point"`
	TotalTask  int      `firestore:"total_task"`
}

func (sc SubCategory) EncID() string {
	return aes.EncryptID(sc.ID)
}
