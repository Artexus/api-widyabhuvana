package subtask

import (
	"context"

	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/subtask"

	"cloud.google.com/go/firestore"
)

type Repository struct {
	client *firestore.Client
}

func NewRepository(client *firestore.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r Repository) subtask() *firestore.CollectionRef {
	return r.client.Collection("sub_tasks")
}

func (r Repository) Get(ctx context.Context, id string) (entity db.SubTask, err error) {
	snap, err := r.subtask().Doc(id).Get(ctx)
	if err != nil {
		return
	}

	snap.DataTo(&entity)
	return
}
