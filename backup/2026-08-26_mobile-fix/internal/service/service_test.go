package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

// ─────────── تست رمزنگاری AES-256-GCM ───────────

func TestCryptoEncryptDecrypt(t *testing.T) {
	c := NewCryptoService("test-secret-key-1234567890")

	plaintext := []byte("سلام خان! این یک پیام محرمانه است")
	ciphertext, err := c.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if string(ciphertext) == string(plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := c.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Fatalf("roundtrip mismatch: got %q want %q", decrypted, plaintext)
	}
	t.Log("✅ AES-256-GCM roundtrip OK")
}

func TestCryptoWrongKey(t *testing.T) {
	c1 := NewCryptoService("secret-key-one-123456")
	c2 := NewCryptoService("secret-key-two-123456")

	ciphertext, err := c1.Encrypt([]byte("data"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c2.Decrypt(ciphertext); err == nil {
		t.Fatal("decrypt with wrong key should fail")
	}
	t.Log("✅ wrong key correctly rejected")
}

// ─────────── تست لایسنس Ed25519 ───────────

func signLicense(t *testing.T, priv ed25519.PrivateKey, payload map[string]interface{}) []byte {
	t.Helper()
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	sig := ed25519.Sign(priv, data)
	return append(data, []byte("\n"+base64.StdEncoding.EncodeToString(sig))...)
}

func splitLicense(t *testing.T, lic []byte) (data, sig []byte) {
	t.Helper()
	for i := len(lic) - 1; i >= 0; i-- {
		if lic[i] == '\n' {
			data = lic[:i]
			sig, err := base64.StdEncoding.DecodeString(string(lic[i+1:]))
			if err != nil {
				t.Fatal(err)
			}
			return data, sig
		}
	}
	t.Fatal("invalid license format")
	return nil, nil
}

func TestLicenseValid(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	payload := map[string]interface{}{
		"v":   1,
		"t":   "شرکت تست",
		"m_u": 100,
		"i":   time.Now().UTC().Format(time.RFC3339),
		"e":   "2027-12-31T00:00:00Z",
		"f":   []string{"chat", "files", "reactions"},
	}

	lic := signLicense(t, priv, payload)

	data, sig := splitLicense(t, lic)
	if !ed25519.Verify(pub, data, sig) {
		t.Fatal("valid license signature rejected")
	}
	t.Log("✅ valid license verified")
}

func TestLicenseTampered(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	payload := map[string]interface{}{
		"v":   1,
		"t":   "شرکت تست",
		"m_u": 100, // امضا شده با ۱۰۰
		"i":   time.Now().UTC().Format(time.RFC3339),
		"e":   "2027-12-31T00:00:00Z",
	}

	lic := signLicense(t, priv, payload)

	// دستکاری: ۱۰۰ → ۹۹۹
	data, sig := splitLicense(t, lic)
	var tampered map[string]interface{}
	if err := json.Unmarshal(data, &tampered); err != nil {
		t.Fatal(err)
	}
	tampered["m_u"] = 999
	newData, _ := json.Marshal(tampered)

	// تأیید باید شکست بخورد
	if ed25519.Verify(priv.Public().(ed25519.PublicKey), newData, sig) {
		t.Fatal("tampered license should fail verification")
	}
	t.Log("✅ tampered license correctly rejected")
}
