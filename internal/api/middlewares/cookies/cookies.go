package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"io"
	"net/http"
	"strings"
)

const Title = "Authorization"

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
	ErrTokenBlank   = errors.New("missing token in Authorization header and Cookie")
)

func Write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	// Check the total length of the cookie contents. Return the ErrValueTooLong
	// error if it's more than 4096 bytes.
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	w.Header().Set(cookie.Name, cookie.Value)
	http.SetCookie(w, &cookie)

	return nil
}

func Read(r *http.Request, name string) (string, error) {
	token, err := fetchToken(r, name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", ErrInvalidValue
	}
	// Return the decoded cookie value.
	return string(value), nil
}

func WriteEncrypted(w http.ResponseWriter, cookie http.Cookie) error {
	block, err := aes.NewCipher(config.CookieSalt)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = string(encryptedValue)

	// Write the cookie as normal.
	return Write(w, cookie)
}

func ReadEncrypted(r *http.Request, name string) (string, error) {
	// Read the encrypted value from the cookie as normal.
	encryptedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(config.CookieSalt)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", ErrInvalidValue
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", ErrInvalidValue
	}

	// Check that the cookie name is the expected one and hasn't been changed.
	if expectedName != name {
		return "", ErrInvalidValue
	}

	// Return the plaintext cookie value.
	return value, nil
}

func fetchToken(r *http.Request, name string) (string, error) {
	authHeader := r.Header.Get(name)
	if authHeader != "" {
		return authHeader, nil
	}

	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", ErrTokenBlank
	}

	return cookie.Value, nil
}
