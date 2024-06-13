package category

import "github.com/Artexus/api-widyabhuvana/src/util/aes"

type Category struct {
	ID            string   `firestore:"id"`
	Name          string   `firestore:"name"`
	Description   string   `firestore:"description"`
	SubCategories []string `firestore:"sub_categories"`
}

func (c Category) EncID() string {
	return aes.EncryptID(c.ID)
}
