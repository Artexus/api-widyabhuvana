package constant

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load()

	JWTInterval, _ = time.ParseDuration(os.Getenv("JWT_INTERVAL"))

	JWTSignedKey = os.Getenv("TOKEN_SIGNED_KEY")

	PasswordMaxlength = 8

	Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
)
