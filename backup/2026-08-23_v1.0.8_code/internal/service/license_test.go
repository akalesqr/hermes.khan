package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// ─────────── تست یکپارچه لایسنس ───────────

func TestLicenseIntegration(t *testing.T) {
	// ۱. کلیدها
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	SetPublicKey(pub)

	// ۲. حالت اولیه: بدون لایسنس → ۲۰ کاربر
	if got := LicenseMaxUsers(); got != 20 {
		t.Fatalf("no license should allow 20 users, got %d", got)
	}
	t.Log("✅ no license → 20 users")

	// ۳. ساخت لایسنس معتبر ۱۰۰ کاربری
	dir := t.TempDir()
	licPath := filepath.Join(dir, "license.key")

	payload := LicensePayload{
		Ver:  1,
		To:   "شرکت تست",
		MU:   100,
		Iss:  time.Now().UTC().Format(time.RFC3339),
		Exp:  "2027-12-31T00:00:00Z",
		Feat: []string{"chat", "files"},
	}
	canonical, _ := json.Marshal(struct {
		V  int      `json:"v"`
		T  string   `json:"t"`
		MU int      `json:"m_u"`
		I  string   `json:"i"`
		E  string   `json:"e"`
		F  []string `json:"f"`
	}{payload.Ver, payload.To, payload.MU, payload.Iss, payload.Exp, payload.Feat})
	sig := ed25519.Sign(priv, canonical)
	payload.Sig = base64.StdEncoding.EncodeToString(sig)

	licData, _ := json.Marshal(payload)
	if err := os.WriteFile(licPath, licData, 0600); err != nil {
		t.Fatal(err)
	}

	// ۴. اعمال لایسنس → ۱۰۰ کاربر
	if err := LoadLicense(licPath); err != nil {
		t.Fatalf("LoadLicense: %v", err)
	}
	if got := LicenseMaxUsers(); got != 100 {
		t.Fatalf("valid license should allow 100 users, got %d", got)
	}
	t.Log("✅ valid license → 100 users")

	// ۵. دستکاری → ۵ کاربر (جریمه)
	var tampered LicensePayload
	_ = json.Unmarshal(licData, &tampered)
	tampered.MU = 999
	tamperedData, _ := json.Marshal(tampered)
	if err := os.WriteFile(licPath, tamperedData, 0600); err != nil {
		t.Fatal(err)
	}

	if err := LoadLicense(licPath); err == nil {
		t.Fatal("tampered license should fail to load")
	}
	if got := LicenseMaxUsers(); got != 5 {
		t.Fatalf("tampered license should penalize to 5 users, got %d", got)
	}
	t.Log("✅ tampered license → 5 users (penalty)")

	// ۶. حذف لایسنس → بازگشت به ۲۰
	if err := RemoveLicense(licPath); err != nil {
		t.Fatal(err)
	}
	if got := LicenseMaxUsers(); got != 20 {
		t.Fatalf("remove license should return to 20 users, got %d", got)
	}
	t.Log("✅ remove license → 20 users")
}
