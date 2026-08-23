package service

import (
	"errors"
	"fmt"
	"time"

	"khan/internal/models"
	"khan/internal/repository"
)

// UserService handles user administration (admin/super-admin actions)
type UserService struct {
	users    *repository.UserRepo
	sessions *repository.SessionRepo
	auth     *AuthService
}

func NewUserService(users *repository.UserRepo, sessions *repository.SessionRepo, auth *AuthService) *UserService {
	return &UserService{users: users, sessions: sessions, auth: auth}
}

// CreateUser creates a new user (admin+). Enforces license seat limit.
func (s *UserService) CreateUser(actor *models.User, username, displayName, initialPwd, role string) (*models.User, error) {
	if !actor.CanManageUsers() {
		return nil, errors.New("دسترسی غیرمجاز")
	}

	if err := Validate(initialPwd); err != nil {
		return nil, err
	}

	// License seat check (skip for hidden super admin creation)
	count, err := s.users.CountVisible()
	if err != nil {
		return nil, err
	}
	maxUsers := LicenseMaxUsers()
	if count >= maxUsers {
		return nil, fmt.Errorf("ظرفیت مجاز (%d کاربر) تکمیل شده است", maxUsers)
	}

	existing, err := s.users.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("این نام کاربری قبلاً ثبت شده است")
	}

	// Role validation: only super admin can create admins
	if role == models.RoleAdmin && !actor.CanManageAdmins() {
		return nil, errors.New("فقط مدیر اصلی می‌تواند ادمین بسازد")
	}
	if role == models.RoleSuperAdmin {
		return nil, errors.New("نقش غیرمجاز")
	}
	if role == "" {
		role = models.RoleUser
	}

	hash, err := s.auth.pass.Hash(initialPwd)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username:      username,
		PasswordHash:  hash,
		DisplayName:   displayName,
		Role:          role,
		Active:        true,
		MustChangePwd: true, // force first login password change
		CreatedAt:     time.Now(),
	}
	id, err := s.users.Create(u)
	if err != nil {
		return nil, err
	}
	u.ID = id
	return u.Sanitized(), nil
}

// DeleteUser removes a user (admin+)
func (s *UserService) DeleteUser(actor *models.User, userID int64) error {
	if !actor.CanManageUsers() {
		return errors.New("دسترسی غیرمجاز")
	}
	target, err := s.users.GetByID(userID)
	if err != nil || target == nil {
		return errors.New("کاربر یافت نشد")
	}
	// Can't delete admins unless super admin; can't delete super admin ever
	if target.Role == models.RoleAdmin && !actor.CanManageAdmins() {
		return errors.New("فقط مدیر اصلی می‌تواند ادمین را حذف کند")
	}
	if target.Role == models.RoleSuperAdmin {
		return errors.New("این کاربر قابل حذف نیست")
	}
	if target.ID == actor.ID {
		return errors.New("نمی‌توانید خودتان را حذف کنید")
	}
	return s.users.Delete(userID)
}

// ResetPassword generates a temp password and forces change (admin+)
func (s *UserService) ResetPassword(actor *models.User, userID int64, newPwd string) error {
	if !actor.CanManageUsers() {
		return errors.New("دسترسی غیرمجاز")
	}
	target, err := s.users.GetByID(userID)
	if err != nil || target == nil {
		return errors.New("کاربر یافت نشد")
	}
	if target.Role == models.RoleSuperAdmin {
		return errors.New("این کاربر قابل مدیریت نیست")
	}
	_ = s.sessions.DeleteForUser(userID) // force re-login
	return s.auth.SetPassword(userID, newPwd, true)
}

// ToggleActive suspends/reactivates a user (admin+)
func (s *UserService) ToggleActive(actor *models.User, userID int64) error {
	if !actor.CanManageUsers() {
		return errors.New("دسترسی غیرمجاز")
	}
	target, err := s.users.GetByID(userID)
	if err != nil || target == nil {
		return errors.New("کاربر یافت نشد")
	}
	if target.Role == models.RoleSuperAdmin {
		return errors.New("این کاربر قابل مدیریت نیست")
	}
	if target.Role == models.RoleAdmin && !actor.CanManageAdmins() {
		return errors.New("فقط مدیر اصلی می‌تواند ادمین را مدیریت کند")
	}
	target.Active = !target.Active
	if !target.Active {
		_ = s.sessions.DeleteForUser(userID)
	}
	return s.users.Update(target)
}

// SetRole changes a user's role (admin+ for supervisor, super-admin for admin)
func (s *UserService) SetRole(actor *models.User, userID int64, newRole string) error {
	target, err := s.users.GetByID(userID)
	if err != nil || target == nil {
		return errors.New("کاربر یافت نشد")
	}
	if target.Role == models.RoleSuperAdmin {
		return errors.New("این کاربر قابل تغییر نیست")
	}

	switch newRole {
	case models.RoleUser, models.RoleSupervisor:
		if !actor.CanManageUsers() {
			return errors.New("دسترسی غیرمجاز")
		}
	case models.RoleAdmin:
		if !actor.CanManageAdmins() {
			return errors.New("فقط مدیر اصلی می‌تواند ادمین بسازد")
		}
	case models.RoleSuperAdmin:
		return errors.New("نقش غیرمجاز")
	default:
		return errors.New("نقش نامعتبر")
	}

	// Can't demote an admin unless super admin
	if target.Role == models.RoleAdmin && newRole != models.RoleAdmin && !actor.CanManageAdmins() {
		return errors.New("فقط مدیر اصلی می‌تواند ادمین را تنزل دهد")
	}

	target.Role = newRole
	return s.users.Update(target)
}
