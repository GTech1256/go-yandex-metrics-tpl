package hashguard

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/logging"
	"io"
	"net/http"
)

const HeaderHash = "HashSHA256"

func getEncryptBodyHash(body, key []byte) ([]byte, error) {
	// подписываем алгоритмом HMAC, используя SHA-256
	hash := hmac.New(sha256.New, key)
	hash.Write(body)
	dst := hash.Sum(nil)

	return dst, nil

}

func decodeHash(hash string) ([]byte, error) {
	return hex.DecodeString(hash)
}

func validateHash(hash string, key, body []byte) (bool, error) {
	hashByte, err := decodeHash(hash)
	if err != nil {
		return false, err
	}

	testHash, err := getEncryptBodyHash(body, key)
	if err != nil {
		return false, err
	}

	return hmac.Equal(testHash, hashByte), nil
}

func WithHashGuard(h http.Handler, key []byte, logger logging.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		hashHeader := r.Header.Get(HeaderHash)

		if hashHeader != "" {
			logger.Infof("Валидация хэша из Header ", HeaderHash)
			bodyByte, err := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
			if err != nil {
				logger.Error(err)
				return
			}
			isValidHash, err := validateHash(hashHeader, key, bodyByte)
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !isValidHash {
				logger.Infof("Хэш невалиден")
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				logger.Infof("Хэш валиден")
			}
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
