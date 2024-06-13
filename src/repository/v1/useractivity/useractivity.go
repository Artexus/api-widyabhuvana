package useractivity

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Artexus/api-widyabhuvana/src/constant"
	db "github.com/Artexus/api-widyabhuvana/src/entity/v1/db/useractivity"
	"github.com/gin-gonic/gin"
)

type Repository struct {
	client *firestore.Client
}

func NewRepository(client *firestore.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r Repository) useractivity() *firestore.CollectionRef {
	return r.client.Collection("user_activities")
}

func (r Repository) Get(ctx context.Context, userID, categoryID string) (docID string, entity db.UserActivity, err error) {
	it := r.useractivity().
		WhereEntity(firestore.AndFilter{
			Filters: []firestore.EntityFilter{
				firestore.PropertyFilter{
					Path:     "user_id",
					Operator: "==",
					Value:    userID,
				},
				firestore.PropertyFilter{
					Path:     "category_id",
					Operator: "==",
					Value:    categoryID,
				},
			},
		}).Limit(1).Documents(ctx)

	snap, err := it.GetAll()
	if len(snap) == 0 {
		err = constant.ErrNotFound
		return
	}

	if err != nil {
		return
	}

	docID = snap[0].Ref.ID
	snap[0].DataTo(&entity)
	return
}

func (r Repository) GetByCategoryIDs(ctx context.Context, userID string, categoryIDs []string) (entities []db.UserActivity, err error) {
	if len(categoryIDs) == 0 {
		return
	}

	it := r.useractivity().
		WhereEntity(firestore.AndFilter{
			Filters: []firestore.EntityFilter{
				firestore.PropertyFilter{
					Path:     "user_id",
					Operator: "==",
					Value:    userID,
				},
				firestore.PropertyFilter{
					Path:     "category_id",
					Operator: "in",
					Value:    categoryIDs,
				},
			},
		}).Limit(1).Documents(ctx)

	snaps, err := it.GetAll()
	if len(snaps) == 0 {
		err = constant.ErrNotFound
		return
	}

	if err != nil {
		return
	}

	for _, snap := range snaps {
		entity := db.UserActivity{}

		snap.DataTo(&entity)
		entities = append(entities, entity)
	}
	return
}

func (r Repository) GetByUserID(ctx context.Context, userID string) (entities []db.UserActivity, err error) {
	it := r.useractivity().
		WhereEntity(firestore.AndFilter{
			Filters: []firestore.EntityFilter{
				firestore.PropertyFilter{
					Path:     "user_id",
					Operator: "==",
					Value:    userID,
				},
				firestore.PropertyFilter{
					Path:     "status",
					Operator: "==",
					Value:    db.NotYet,
				},
			},
		}).Limit(1).Documents(ctx)

	snaps, err := it.GetAll()
	if len(snaps) == 0 {
		err = constant.ErrNotFound
		return
	}

	if err != nil {
		return
	}

	for _, snap := range snaps {
		entity := db.UserActivity{}

		snap.DataTo(&entity)
		entities = append(entities, entity)
	}

	return
}

func (r Repository) Set(ctx *gin.Context, entity db.UserActivity, docID string) (err error) {
	q := r.useractivity()
	if docID != "" {
		_, err = q.Doc(docID).
			Set(ctx, entity)
	} else {
		_, _, err = q.Add(ctx, entity)
	}

	return
}
