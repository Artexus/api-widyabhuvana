package user

import (
	"github.com/Artexus/api-widyabhuvana/src/util/aes"
)

type User struct {
	ID         string `firestore:"id"`
	Email      string `firestore:"email"`
	Username   string `firestore:"username"`
	TotalPoint int    `firestore:"total_point"`
	Password   string `firestore:"password"`
	PhotoURL   string `firestore:"photo_url"`
	CreatedAt  int64  `firestore:"created_at"`
	DeletedAt  int64  `firestore:"deleted_at"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) EncID() string {
	return aes.EncryptID(u.ID)
}
