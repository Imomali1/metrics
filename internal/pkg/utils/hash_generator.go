package utils

import (
	"crypto/hmac"
	"crypto/sha256"
)

func GenerateHash(value []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(value)
	return string(h.Sum(nil))
}
