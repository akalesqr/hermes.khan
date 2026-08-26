package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// CryptoService handles AES-256-GCM encryption of messages at rest
type CryptoService struct {
	key []byte
}

// NewCryptoService derives an AES key from the configured secret
func NewCryptoService(secret string) *CryptoService {
	sum := sha256.Sum256([]byte(secret))
	return &CryptoService{key: sum[:]}
}

// Encrypt encrypts plaintext with AES-256-GCM. Returns base64(iv).base64(ciphertext)
func (c *CryptoService) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, iv, plaintext, nil)

	// iv || ciphertext (with tag appended by GCM)
	out := make([]byte, 0, len(iv)+len(ciphertext))
	out = append(out, iv...)
	out = append(out, ciphertext...)
	return []byte(base64.StdEncoding.EncodeToString(out)), nil
}

// Decrypt reverses Encrypt
func (c *CryptoService) Decrypt(data []byte) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(raw) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	iv, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	return gcm.Open(nil, iv, ciphertext, nil)
}

// EncryptString is a convenience wrapper
func (c *CryptoService) EncryptString(s string) ([]byte, error) {
	return c.Encrypt([]byte(s))
}

// DecryptString is a convenience wrapper
func (c *CryptoService) DecryptString(b []byte) (string, error) {
	plain, err := c.Decrypt(b)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// GenerateKey produces a random 32-byte base64 key for config
func GenerateKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
