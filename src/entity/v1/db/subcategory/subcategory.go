package subcategory

import "cloud.google.com/go/firestore"

type SubCategory struct {
	CategoryID string                    `firestore:"category_id"`
	Name       string                    `firestore:"name"`
	Tasks      []firestore.CollectionRef `firestore:"tasks"`
	MaxPoint   int                       `firestore:"max_point"`
}
