package aes

import (
	"encoding/hex"
)

func EncryptID(id string) string {
	return hex.EncodeToString([]byte(id))
}

func DecryptID(encID string) (res string, err error) {
	resByte, err := hex.DecodeString(encID)
	if err != nil {
		return
	}

	res = string(resByte)
	return
}
