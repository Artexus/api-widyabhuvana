package category

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/category"
	"github.com/Artexus/api-widyabhuvana/src/util/pagination"
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

func (r Repository) category() *firestore.CollectionRef {
	return r.client.Collection("categories")
}

func (r Repository) Get(ctx context.Context, pgn pagination.Pagination) (entities []db.Category, err error) {
	entities = make([]db.Category, 0)
	snaps, err := r.category().
		Limit(pgn.Limit).
		Offset(pgn.Offset).
		Documents(ctx).GetAll()
	if err != nil {
		return
	}

	for _, ss := range snaps {
		entity := db.Category{}
		ss.DataTo(&entity)

		entities = append(entities, entity)
	}

	return
}

func (r Repository) GetByIDs(ctx context.Context, ids []string) (entities []db.Category, err error) {
	if len(ids) == 0 {
		return
	}

	entities = make([]db.Category, 0)
	snaps, err := r.category().
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
		entity := db.Category{}
		ss.DataTo(&entity)

		entities = append(entities, entity)
	}

	return
}

func (r Repository) Take(ctx context.Context, id string) (entity db.Category, err error) {
	snap, err := r.category().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		err = constant.ErrNotFound
		return
	} else if err != nil {
		return
	}

	err = snap.DataTo(&entity)
	return
}
