package category

import (
	"context"

	"cloud.google.com/go/firestore"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/category"
	"github.com/Artexus/api-widyabhuvana/src/util/pagination"
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
