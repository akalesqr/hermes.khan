package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordService hashes and verifies passwords with Argon2id
type PasswordService struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// NewPasswordService returns an Argon2id hasher with strong defaults
func NewPasswordService() *PasswordService {
	return &PasswordService{
		time:    3,
		memory:  64 * 1024, // 64MB
		threads: 1,
		keyLen:  32,
	}
}

// Hash creates an Argon2id hash: "$argon2id$v=19$m=65536,t=3,p=1$salt$hash"
func (p *PasswordService) Hash(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threads, p.keyLen)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.memory, p.time, p.threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash)), nil
}

// Verify checks a password against a stored hash
func (p *PasswordService) Verify(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var memory uint32
	var time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expected)))

	if subtle.ConstantTimeCompare(hash, expected) == 1 {
		return true, nil
	}
	return false, nil
}

// Validate enforces password policy (min 8 chars, mixed)
func Validate(password string) error {
	if len(password) < 8 {
		return errors.New("رمز عبور باید حداقل ۸ کاراکتر باشد")
	}
	hasUpper, hasLower, hasDigit, hasSymbol := false, false, false, false
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;':,./<>?", r):
			hasSymbol = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("رمز عبور باید شامل حروف بزرگ، حروف کوچک و عدد باشد")
	}
	if !hasSymbol && len(password) < 12 {
		return errors.New("رمز عبور ضعیف است")
	}
	_ = hex.EncodeToString // keep import
	return nil
}
