package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(value []byte, key string) string {
	h := sha256.New()
	// передаём байты с ключом для хеширования
	value = append(value, []byte(key)...)
	h.Write(value)
	// вычисляем хеш
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}
