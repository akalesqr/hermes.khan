package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"khan/config"
	"khan/internal/models"
	"khan/internal/repository"
)

// AuthService handles login/logout/token management
type AuthService struct {
	users    *repository.UserRepo
	sessions *repository.SessionRepo
	pass     *PasswordService
	cfg      *config.Config
}

func NewAuthService(users *repository.UserRepo, sessions *repository.SessionRepo, cfg *config.Config) *AuthService {
	return &AuthService{
		users:    users,
		sessions: sessions,
		pass:     NewPasswordService(),
		cfg:      cfg,
	}
}

// Login validates credentials and issues a session token
func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	u, err := s.users.GetByUsername(username)
	if err != nil {
		return nil, "", err
	}
	if u == nil {
		return nil, "", errors.New("نام کاربری یا رمز عبور اشتباه است")
	}
	if !u.Active {
		return nil, "", errors.New("حساب کاربری غیرفعال است")
	}

	ok, err := s.pass.Verify(password, u.PasswordHash)
	if err != nil || !ok {
		return nil, "", errors.New("نام کاربری یا رمز عبور اشتباه است")
	}

	token, err := s.newToken()
	if err != nil {
		return nil, "", err
	}

	sess := &models.Session{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.Security.JWTExpireDays) * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
	if err := s.sessions.Create(sess); err != nil {
		return nil, "", err
	}

	_ = s.users.TouchLastSeen(u.ID)
	return u, token, nil
}

// Logout invalidates a session
func (s *AuthService) Logout(token string) error {
	return s.sessions.Delete(token)
}

// ValidateToken checks a session token and returns the user
func (s *AuthService) ValidateToken(token string) (*models.User, error) {
	sess, err := s.sessions.GetByToken(token)
	if err != nil {
		return nil, errors.New("نشست نامعتبر است")
	}
	if sess.ExpiresAt.Before(time.Now()) {
		_ = s.sessions.Delete(token)
		return nil, errors.New("نشست منقضی شده است")
	}
	u, err := s.users.GetByID(sess.UserID)
	if err != nil || u == nil {
		return nil, errors.New("کاربر یافت نشد")
	}
	if !u.Active {
		return nil, errors.New("حساب کاربری غیرفعال است")
	}
	return u, nil
}

// ChangePassword updates the user's password (requires current password)
func (s *AuthService) ChangePassword(userID int64, currentPwd, newPwd string) error {
	u, err := s.users.GetByID(userID)
	if err != nil || u == nil {
		return errors.New("کاربر یافت نشد")
	}
	ok, err := s.pass.Verify(currentPwd, u.PasswordHash)
	if err != nil || !ok {
		return errors.New("رمز عبور فعلی اشتباه است")
	}
	if err := Validate(newPwd); err != nil {
		return err
	}
	hash, err := s.pass.Hash(newPwd)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(userID, hash)
}

// SetPassword sets a new password (admin reset / first setup)
func (s *AuthService) SetPassword(userID int64, newPwd string, forceChange bool) error {
	if err := Validate(newPwd); err != nil {
		return err
	}
	hash, err := s.pass.Hash(newPwd)
	if err != nil {
		return err
	}
	u, err := s.users.GetByID(userID)
	if err != nil || u == nil {
		return errors.New("کاربر یافت نشد")
	}
	u.PasswordHash = hash
	u.MustChangePwd = forceChange
	return s.users.Update(u)
}

func (s *AuthService) newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// PassHash hashes a plaintext password (used by setup/bootstrap)
func (s *AuthService) PassHash(password string) (string, error) {
	return s.pass.Hash(password)
}
