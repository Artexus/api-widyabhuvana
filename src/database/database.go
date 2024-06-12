package database

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func InitDB() (client *firestore.Client, err error) {
	_ = godotenv.Load()
	client, err = firestore.NewClient(context.Background(), firestore.DetectProjectID, option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_API_KEY"))))
	return
}
