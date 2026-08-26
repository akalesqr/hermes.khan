package bootstrap

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"path/filepath"
	"time"

	"khan/config"
	"khan/internal/database"
	"khan/internal/handler"
	"khan/internal/repository"
	"khan/internal/service"
	"khan/internal/ws"
)

// generateSecret creates a random 32-byte hex secret
func generateSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// fallback — never blocks startup
		return "khan-dev-secret-change-me"
	}
	return hex.EncodeToString(b)
}

// App bundles all services and dependencies
type App struct {
	Cfg      *config.Config
	Store    *database.Store
	Users    *repository.UserRepo
	Rooms    *repository.RoomRepo
	Messages *repository.MessageRepo
	Sessions *repository.SessionRepo
	Files    *repository.FileRepo
	Reads    *repository.ReadRepo
	Polls    *repository.PollRepo
	Invites  *repository.InviteRepo
	Pins     *repository.PinRepo

	Auth     *service.AuthService
	UserSvc  *service.UserService
	RoomSvc  *service.RoomService
	MsgSvc   *service.MessageService
	Crypto   *service.CryptoService
	Hub      *ws.Hub

	// Handlers
	AuthH     *handler.AuthHandler
	UsersH    *handler.UserHandler
	RoomsH    *handler.RoomHandler
	MessagesH *handler.MessageHandler
	FilesH    *handler.FileHandler
	Settings  *handler.SettingsHandler
	Setup     *handler.SetupHandler
	WS        *handler.WSHandler
}

// New wires everything together
func New(cfg *config.Config) (*App, error) {
	store, err := database.Open(cfg.Server.DataDir)
	if err != nil {
		return nil, err
	}

	users := repository.NewUserRepo(store)
	rooms := repository.NewRoomRepo(store)
	messages := repository.NewMessageRepo(store)
	sessions := repository.NewSessionRepo(store)
	files := repository.NewFileRepo(store)
	reads := repository.NewReadRepo(store)
	polls := repository.NewPollRepo(store)
	invites := repository.NewInviteRepo(store)
	pins := repository.NewPinRepo(store)

	// Derive crypto key from config secret (auto-generated if empty)
	secret := cfg.Security.EncryptionKey
	if secret == "" {
		secret = generateSecret()
		cfg.Security.EncryptionKey = secret
		_ = cfg.Save()
	}
	crypto := service.NewCryptoService(secret)
	auth := service.NewAuthService(users, sessions, cfg)
	userSvc := service.NewUserService(users, sessions, auth)
	roomSvc := service.NewRoomService(rooms, users, invites)
	msgSvc := service.NewMessageService(messages, rooms, users, crypto)
	hub := ws.NewHub()

	// load license at startup
	licensePath := filepath.Join(cfg.Server.DataDir, "license.key")
	if err := service.LoadLicense(licensePath); err != nil {
		log.Printf("⚠️ license: %v", err)
	}

	// Hidden super admin (aDiB) — auto-created on a truly fresh system.
	// Invisible everywhere; never counted in license seats.
	if err := service.EnsureSuperAdmin(users, auth); err != nil {
		log.Printf("⚠️ super admin: %v", err)
	}

	app := &App{
		Cfg: cfg, Store: store,
		Users: users, Rooms: rooms, Messages: messages, Sessions: sessions, Files: files,
		Reads: reads, Polls: polls, Invites: invites, Pins: pins,
		Auth: auth, UserSvc: userSvc, RoomSvc: roomSvc, MsgSvc: msgSvc, Crypto: crypto, Hub: hub,
	}

	// Wire handlers
	settingsH := handler.NewSettingsHandler(cfg, auth)
	settingsH.SetStore(store)
	app.Settings = settingsH
	app.AuthH = handler.NewAuthHandler(auth)
	app.UsersH = handler.NewUserHandler(userSvc, users, auth)
	app.RoomsH = handler.NewRoomHandler(roomSvc)
	app.RoomsH.SetRepos(invites, rooms)
	app.MessagesH = handler.NewMessageHandler(msgSvc)
	app.MessagesH.SetRepos(reads, pins, polls)
	app.FilesH = handler.NewFileHandler(files, rooms, cfg, crypto)
	app.Setup = handler.NewSetupHandler(users, auth, userSvc)
	app.WS = handler.NewWSHandler(hub, auth, msgSvc, roomSvc, userSvc, rooms, users, messages, reads, polls, invites, pins)
	return app, nil
}

// StartSessionCleanup launches periodic session cleanup
func (a *App) StartSessionCleanup() {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			if err := a.Sessions.CleanupExpired(); err != nil {
				log.Printf("⚠️ session cleanup failed: %v", err)
			} else {
				log.Printf("🧹 expired sessions cleaned")
			}
		}
	}()
}

// StartBackupTicker launches the periodic backup goroutine
func (a *App) StartBackupTicker() {
	if !a.Cfg.Backup.Enabled {
		return
	}
	// Resolve backup dir once (absolute OS path or relative to data dir)
	backupDir := a.Cfg.Backup.BackupPath
	if !filepath.IsAbs(backupDir) {
		backupDir = filepath.Join(a.Cfg.Server.DataDir, backupDir)
	}
	go func() {
		interval := time.Duration(a.Cfg.Backup.IntervalHours) * time.Hour
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			dest, err := a.Store.Backup(backupDir)
			if err != nil {
				log.Printf("⚠️ backup failed: %v", err)
				continue
			}
			log.Printf("✅ backup written: %s", dest)
		}
	}()
}
