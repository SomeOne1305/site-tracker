package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"visit-tracker/config"
)

// =============================
// REFRESH TOKEN GENERATION
// =============================
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// =============================
// LOAD & VALIDATE KEY
// =============================
func loadKey() ([]byte, error) {
	cfg := config.LoadConfig()
	key := []byte(cfg.EncryptionKey)

	if len(key) != 32 {
		return nil, errors.New("encryption key must be exactly 32 bytes (AES-256)")
	}

	return key, nil
}

// =============================
// ENCRYPT STRING (FIXED)
// =============================
func EncryptID(id string) (string, error) {
	key, err := loadKey()
	if err != nil {
		fmt.Println("Key err: ", err)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Block err: ", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {

		fmt.Println("GCM err: ", err)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("IO err: ", err)
		return "", err
	}

	// ✅ FIX: no fmt.Sprintf, use raw string
	plaintext := []byte(id)

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// =============================
// DECRYPT STRING
// =============================
func DecryptID(tokenB64 string) (string, error) {
	key, err := loadKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(tokenB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("malformed token")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err // tampered or wrong key
	}

	return string(plaintext), nil
}
