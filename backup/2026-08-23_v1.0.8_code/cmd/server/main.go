package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"khan/config"
	"khan/internal/bootstrap"
	"khan/internal/handler"
	"khan/internal/service"
)

var version = "1.0.8"

//go:embed all:web
var webFS embed.FS
func main() {
	var (
		configPath  = flag.String("config", "config.json", "path to config file")
		dataDir     = flag.String("data", "data", "data directory")
		port        = flag.Int("port", 0, "override port (1700-1799)")
		versionFlag = flag.Bool("version", false, "print version")
	)
	flag.Parse()

	if *versionFlag {
		log.Printf("Khan v%s", version)
		return
	}

	// Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	if *dataDir != "" {
		cfg.Server.DataDir = *dataDir
		cfg.Server.DBPath = filepath.Join(*dataDir, "khan.db.json")
	}
	if *port != 0 {
		if *port < 1700 || *port > 1799 {
			log.Fatalf("port must be 1700-1799 (got %d)", *port)
		}
		cfg.Server.Port = *port
	}
	// PORT env override (Replit, Render, Fly.io, etc.)
	// — allows any port 1-65535, bypasses the 1700-1799 LAN range
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil && p > 0 && p <= 65535 {
			cfg.Server.Port = p
			log.Printf("🌍 PORT env: listening on :%d", p)
		}
	}
	config.SetInstance(cfg)

	// Bootstrap app
	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatalf("bootstrap error: %v", err)
	}
	app.StartBackupTicker()
	app.StartSessionCleanup()

	// Build router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ---- Public endpoints ----
	r.Get("/api/settings/info", app.Settings.Info)
	r.Get("/api/setup/needs-setup", app.Setup.NeedsSetup)
	r.Post("/panel/setup", app.Setup.Setup)
	r.Post("/api/auth/login", app.AuthH.Login)

	// ---- Auth required ----
	r.Group(func(r chi.Router) {
		r.Use(handler.RequireAuth(app.Auth))
		r.Post("/api/auth/logout", app.AuthH.Logout)
		r.Get("/api/auth/me", app.AuthH.Me)
		r.Post("/api/auth/change-password", app.AuthH.ChangePassword)

		r.Post("/api/users/avatar", app.UsersH.UploadAvatar)
		r.Put("/api/users/me", app.UsersH.UpdateProfile)

		r.Get("/api/users", app.UsersH.List)
		r.Get("/api/users/export", app.UsersH.Export)
		r.Get("/api/users/search", app.UsersH.Search)
		r.Post("/api/users", app.UsersH.Create)
		r.Delete("/api/users/{id}", app.UsersH.Delete)
		r.Post("/api/users/{id}/reset-password", app.UsersH.ResetPassword)
		r.Post("/api/users/{id}/toggle-active", app.UsersH.ToggleActive)
		r.Post("/api/users/{id}/role", app.UsersH.SetRole)

		r.Get("/api/rooms", app.RoomsH.List)
		r.Post("/api/rooms", app.RoomsH.Create)
		r.Get("/api/rooms/public", app.RoomsH.PublicRooms)
		r.Get("/api/rooms/private", app.RoomsH.PrivateRooms)
		r.Post("/api/rooms/{id}/join", app.RoomsH.Join)
		r.Post("/api/rooms/{id}/invite", app.RoomsH.Invite)
		r.Post("/api/rooms/{id}/members", app.RoomsH.AddMember)
		r.Delete("/api/rooms/{id}/members/{uid}", app.RoomsH.RemoveMember)
		r.Get("/api/rooms/{id}/members", app.RoomsH.Members)
		r.Post("/api/rooms/{id}/rename", app.RoomsH.Rename)
		r.Post("/api/rooms/dm/{uid}", app.RoomsH.StartDM)

		r.Get("/api/departments", app.RoomsH.Departments)
		r.Post("/api/departments", app.RoomsH.CreateDepartment)
		r.Delete("/api/departments/{id}", app.RoomsH.DeleteDepartment)

		r.Get("/api/messages/{id}", app.MessagesH.List)
		r.Get("/api/messages/search", app.MessagesH.Search)
		r.Get("/api/messages/{id}/pins", app.MessagesH.Pins)
		r.Post("/api/messages/{id}/edit", app.MessagesH.Edit)
		r.Delete("/api/messages/{id}", app.MessagesH.Delete)
		r.Post("/api/messages/{id}/forward", app.MessagesH.Forward)
		r.Post("/api/messages/{id}/urgent", app.MessagesH.ToggleUrgent)
		r.Post("/api/messages/{id}/reactions", app.MessagesH.AddReaction)
		r.Delete("/api/messages/{id}/reactions", app.MessagesH.RemoveReaction)

		r.Post("/api/rooms/{id}/files", app.FilesH.Upload)
		r.Get("/api/files/{id}/download", app.FilesH.Download)

		})

		// ---- Admin only ----
		r.Group(func(r chi.Router) {
			r.Use(handler.RequireAuth(app.Auth))
			r.Use(handler.RequireRole("admin"))

			r.Get("/api/settings/license", app.Settings.LicenseStatus)
			r.Post("/api/settings/license", app.Settings.LicenseApply)
			r.Delete("/api/settings/license", app.Settings.LicenseRemove)
			r.Get("/api/settings/network", app.Settings.NetworkInfo)
			r.Post("/api/settings/network", app.Settings.NetworkUpdate)
			r.Post("/api/settings/backup", app.Settings.Backup)
			r.Get("/api/settings/backups", app.Settings.ListBackups)
			r.Get("/api/settings/backups/{name}/download", app.Settings.DownloadBackup)
			r.Post("/api/settings/restore", app.Settings.RestoreBackup)
			r.Get("/api/settings/logs", app.Settings.Logs)
		})

	// ---- WebSocket (own auth — token via query or header) ----
	r.Get("/ws", app.WS.ServeHTTP)

	// ---- Static files (PWA) ----
	serveStatic(r)

	// ---- Start server (HTTP or HTTPS with optional auto self-signed TLS) ----
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	state, to, errMsg, maxUsers, _ := service.LicenseState()
	scheme := "http"
	if cfg.Server.TLS {
		scheme = "https"
		certPath, keyPath, err := cfg.EnsureSelfSignedCert()
		if err != nil {
			log.Fatalf("TLS cert error: %v", err)
		}
		log.Printf("🔐 TLS enabled — serving HTTPS/WSS (cert=%s)", filepath.Base(certPath)+", key="+filepath.Base(keyPath))
	}
	log.Printf("🏠 Khan v%s listening on %s://%s", version, scheme, addr)
	log.Printf("   License: state=%s max_users=%d to=%q err=%q", state, maxUsers, to, errMsg)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		var err error
		if cfg.Server.TLS {
			certPath, keyPath, cerr := cfg.EnsureSelfSignedCert()
			if cerr != nil {
				log.Fatalf("TLS cert error: %v", cerr)
			}
			err = srv.ListenAndServeTLS(certPath, keyPath)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

// serveStatic serves the embedded PWA with SPA fallback
func serveStatic(r chi.Router) {
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("embed error: %v", err)
	}
	fileServer := http.FileServer(http.FS(sub))

	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/")
		// Try to serve the file
		if path == "" {
			path = "index.html"
		}
		_, err := fs.Stat(sub, path)
		if err != nil {
			// File not found → serve index.html for SPA routing
			req.URL.Path = "/"
			path = "index.html"
		}
		// Add cache headers for static assets
		if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") || 
		   strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".png") ||
		   strings.HasSuffix(path, ".svg") || strings.HasSuffix(path, ".ico") ||
		   strings.HasSuffix(path, ".woff") || strings.HasSuffix(path, ".woff2") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		fileServer.ServeHTTP(w, req)
	}))
}

var _ = fmt.Sprintf
