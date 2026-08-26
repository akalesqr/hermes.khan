package service

import (
	"strings"
	"testing"
)

// ─────────── تست Argon2id ───────────

func TestPasswordHashAndVerify(t *testing.T) {
	p := NewPasswordService()

	password := "Super@Secure123"
	hash, err := p.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	// هش نباید شامل رمز باشد
	if strings.Contains(hash, password) {
		t.Fatal("hash should not contain plaintext password")
	}

	// تأیید درست
	ok, err := p.Verify(password, hash)
	if err != nil {
		t.Fatalf("Verify error: %v", err)
	}
	if !ok {
		t.Fatal("correct password should verify")
	}

	// تأیید غلط
	ok, _ = p.Verify("WrongPassword", hash)
	if ok {
		t.Fatal("wrong password should not verify")
	}
	t.Log("✅ Argon2id hash/verify OK")
}

func TestPasswordDifferentHashes(t *testing.T) {
	p := NewPasswordService()

	h1, _ := p.Hash("same-password")
	h2, _ := p.Hash("same-password")

	// هر هش باید salt متفاوت داشته باشد
	if h1 == h2 {
		t.Fatal("two hashes of same password should differ (random salt)")
	}
	t.Log("✅ unique salts confirmed")
}

func TestPasswordEmpty(t *testing.T) {
	p := NewPasswordService()

	if _, err := p.Hash(""); err == nil {
		t.Log("note: empty password accepted")
	}
}
