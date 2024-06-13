package subcategory

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/subcategory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (r Repository) GetByIDs(ctx context.Context, ids []string) (entities []db.SubCategory, err error) {
	if len(ids) == 0 {
		return
	}

	entities = make([]db.SubCategory, 0)
	snaps, err := r.subcategory().
		WhereEntity(firestore.AndFilter{
			Filters: []firestore.EntityFilter{
				firestore.PropertyFilter{
					Path:     "id",
					Operator: "in",
					Value:    ids,
				},
			},
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

func (r Repository) Take(ctx context.Context, id string) (entity db.SubCategory, err error) {
	snap, err := r.subcategory().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		err = constant.ErrNotFound
		return
	} else if err != nil {
		return
	}

	err = snap.DataTo(&entity)
	return
}
