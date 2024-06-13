package user

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/user"
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

func (r Repository) user() *firestore.CollectionRef {
	return r.client.Collection("users")
}

func (r Repository) Get(ctx context.Context, id string) (entity db.User, err error) {
	snap, err := r.user().Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		err = constant.ErrNotFound
		return
	} else if err != nil {
		return
	}

	err = snap.DataTo(&entity)
	return
}

func (r Repository) TakeByEmail(ctx context.Context, email string) (entity db.User, err error) {
	it := r.user().WhereEntity(firestore.AndFilter{
		Filters: []firestore.EntityFilter{
			firestore.PropertyFilter{Path: "email", Operator: "==", Value: email},
			firestore.PropertyFilter{Path: "deleted_at", Operator: "==", Value: 0},
		},
	}).Limit(1).Documents(ctx)

	docsnaps, err := it.GetAll()
	if err != nil {
		return
	}

	if len(docsnaps) == 0 {
		err = constant.ErrNotFound
		return
	}

	err = docsnaps[0].DataTo(&entity)
	return
}

func (r Repository) Create(ctx context.Context, entity *db.User) (err error) {
	ref := r.user().NewDoc()

	entity.ID = ref.ID
	_, err = ref.Set(ctx, entity)
	if err != nil {
		return
	}

	return
}

func (r Repository) UpdatePoint(ctx context.Context, id string, point int) (err error) {
	_, err = r.user().
		Doc(id).
		Update(ctx, []firestore.Update{{Path: "total_point", Value: firestore.Increment(point)}})
	return
}
