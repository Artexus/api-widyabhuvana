package storage

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

var client *storage.Client

func initStorage() (client *storage.Client, err error) {
	if client == nil {
		var app *firebase.App
		app, err = firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_API_KEY"))))
		if err != nil {
			return
		}

		client, err = app.Storage(context.Background())
	}

	return
}

func Upload(file multipart.File, id, path string) (err error) {
	c, err := initStorage()
	if err != nil {
		return
	}

	bucket, err := c.Bucket(os.Getenv("FIREBASE_BUCKET"))
	if err != nil {
		return
	}

	writer := bucket.Object(path).NewWriter(context.Background())
	defer writer.Close()
	if _, err = io.Copy(writer, file); err != nil {
		return
	}

	return
}
