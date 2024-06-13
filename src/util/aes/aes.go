package aes

import (
	"encoding/hex"

	"github.com/Artexus/api-widyabhuvana/src/constant"
)

func EncryptIDs(ids []string) []string {
	t := []string{}
	for _, id := range ids {
		t = append(t, EncryptID(id))
	}

	return t
}

func EncryptID(id string) string {
	return hex.EncodeToString([]byte(id))
}

func DecryptID(encID string) (res string, err error) {
	if encID == "" {
		err = constant.ErrInvalid
		return
	}

	resByte, err := hex.DecodeString(encID)
	if err != nil {
		return
	}

	res = string(resByte)
	return
}
