package hash

import (
	"crypto/sha256"
	"fmt"
)

func GenerateHashToken(key string) (result string) {

	return fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
}
