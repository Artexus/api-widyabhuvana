package answer

import (
	"context"

	"cloud.google.com/go/firestore"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/answer"
)

type Repository struct {
	client *firestore.Client
}

func NewRepository(client *firestore.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r Repository) answer() *firestore.CollectionRef {
	return r.client.Collection("answers")
}

func (r Repository) Create(ctx context.Context, answer db.Answer) (docID string, err error) {
	q := r.answer().
		NewDoc()

	docID = q.ID
	_, err = q.Create(ctx, answer)
	return
}
