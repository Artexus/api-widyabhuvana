package constant

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type TaskType string

const (
	Learning       TaskType = "LEARNING"
	MultipleChoice TaskType = "MULTIPLE_CHOICE"
	Essay          TaskType = "ESSAY"
	Matching       TaskType = "MATCHING"
	Detective      TaskType = "DETECTIVE"
	Level          TaskType = "LEVEL"
)

var (
	_ = godotenv.Load()

	JWTInterval, _ = time.ParseDuration(os.Getenv("JWT_INTERVAL"))

	JWTSignedKey = os.Getenv("TOKEN_SIGNED_KEY")

	PasswordMaxlength = 8

	Bucket = os.Getenv("FIREBASE_BUCKET")

	Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
)
