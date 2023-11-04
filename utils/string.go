package utils

import (
	"crypto/aes"
	"crypto/cipher"
	cRand "crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/rand"
	"time"
)

var digitRunes []rune = []rune("0123456789")

// RandDigits creates sequence of digit string with expected length.
func RandDigits(length int) string {
	rand.Seed(time.Now().UnixNano())
	res := make([]rune, length)
	for i := range res {
		res[i] = digitRunes[rand.Intn(len(digitRunes))]
	}

	return string(res)
}

// HashString returns hashed string of given string using golang bcrypt algorithm.
func HashString(val string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(val), 10)
	return string(bytes), err
}

// CheckHashedString compare a hashed string with specific plain string.
func CheckHashedString(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

// EncryptText encrypt given string by specific secret
func EncryptText(value, secret string) (string, error) {
	text := []byte(value)
	key := []byte(secret)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("secret is not in correct length")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cRand.Reader, nonce); err != nil {
		return "", err
	}
	res := gcm.Seal(nonce, nonce, text, nil)
	b64 := base64.URLEncoding.EncodeToString(res)
	return b64, nil
}

// DecryptText decrypts given cypher string by specific secret.
func DecryptText(value, secret string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	key := []byte(secret)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", nil
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
