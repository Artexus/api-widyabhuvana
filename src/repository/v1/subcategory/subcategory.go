package subcategory

import (
	"context"

	"cloud.google.com/go/firestore"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/subcategory"
)

type Repository struct {
	client *firestore.Client
}

func NewRepository(client *firestore.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r Repository) subcategory() *firestore.CollectionRef {
	return r.client.Collection("sub_categories")
}

func (r Repository) Get(ctx context.Context, categoryID string) (entities []db.SubCategory, err error) {
	entities = make([]db.SubCategory, 0)
	snaps, err := r.subcategory().
		WhereEntity(firestore.PropertyFilter{
			Path:     "category_id",
			Operator: "==",
			Value:    categoryID,
		}).
		Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, ss := range snaps {
		entity := db.SubCategory{}
		ss.DataTo(&entity)

		entities = append(entities, entity)
	}

	return
}
