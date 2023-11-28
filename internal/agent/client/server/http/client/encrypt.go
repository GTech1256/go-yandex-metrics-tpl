package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (h httpClient) getEncryptBodyHash(body []byte) ([]byte, error) {
	// подписываем алгоритмом HMAC, используя SHA-256
	hash := hmac.New(sha256.New, []byte(*h.HashKey))
	hash.Write(body)
	dst := hash.Sum(nil)

	return dst, nil

}

func (h httpClient) getEncode(src []byte) string {
	return hex.EncodeToString(src)

}
