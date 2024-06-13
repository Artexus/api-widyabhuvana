package task

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/task"
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

func (r Repository) task() *firestore.CollectionRef {
	return r.client.Collection("tasks")
}

func (r Repository) Get(ctx context.Context, id string) (entity db.Task, err error) {
	snap, err := r.task().
		Doc(id).
		Get(ctx)

	if status.Code(err) == codes.NotFound {
		err = constant.ErrNotFound
	}

	if err != nil {
		return
	}

	snap.DataTo(&entity)
	return
}
