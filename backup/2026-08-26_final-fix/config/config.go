package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// getDefaultBackupPath returns OS-appropriate backup directory
func getDefaultBackupPath() string {
	switch runtime.GOOS {
	case "windows":
		// %PROGRAMDATA%\Khan\backups
		programData := os.Getenv("PROGRAMDATA")
		if programData == "" {
			programData = `C:\ProgramData`
		}
		return filepath.Join(programData, "Khan", "backups")
	case "darwin":
		// ~/Library/Application Support/Khan/backups
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "Khan", "backups")
	default:
		// Linux: ~/.local/share/khan/backups
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local", "share", "khan", "backups")
	}
}

// Config represents the server configuration
type Config struct {
	Server       ServerConfig  `json:"server"`
	License      LicenseConfig `json:"license"`
	Security     SecurityConfig  `json:"security"`
	Backup       BackupConfig    `json:"backup"`
	DefaultAdmin DefaultAdmin    `json:"default_admin"`
	path         string          `json:"-"`
}

// SetPath records where this config was loaded from
func (c *Config) SetPath(p string) { c.path = p }

// Save writes config to disk (uses recorded path or explicit path)
func (c *Config) Save(path ...string) error {
	p := c.path
	if len(path) > 0 && path[0] != "" {
		p = path[0]
	}
	if p == "" {
		p = "config.json"
	}
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

type ServerConfig struct {
	Port        int    `json:"port"`
	Host        string `json:"host"`
	AddressType string `json:"address_type"` // "ip" | "dns"
	IP          string `json:"ip"`
	DNS         string `json:"dns"`
	DBPath      string `json:"db_path"`
	DataDir     string `json:"data_dir"`
	MaxUploadMB int    `json:"max_upload_mb"`
	// TLS: when Enabled, the server serves HTTPS (WSS). A self-signed
	// cert is auto-generated into the data dir if Cert/Key are empty.
	TLS      bool   `json:"tls_enabled"`
	CertPath string `json:"tls_cert_path,omitempty"`
	KeyPath  string `json:"tls_key_path,omitempty"`
}

type LicenseConfig struct {
	LicenseFile string `json:"license_file"`
}

type SecurityConfig struct {
	EncryptionKey    string `json:"encryption_key"`
	JWTSecret        string `json:"jwt_secret"`
	JWTExpireDays    int    `json:"jwt_expire_days"`
	MaxLoginAttempts int    `json:"max_login_attempts"`
	LockoutMinutes   int    `json:"lockout_minutes"`
}

type BackupConfig struct {
	Enabled       bool   `json:"enabled"`
	IntervalHours int    `json:"interval_hours"`
	KeepCount     int    `json:"keep_count"`
	BackupPath    string `json:"backup_path"`
}

type DefaultAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Default returns a config with sane defaults
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        1727,
			Host:        "0.0.0.0",
			AddressType: "ip",
			IP:          "192.168.1.100",
			DNS:         "",
			DBPath:      "data/khan.db.json",
			DataDir:     "data",
			MaxUploadMB: 50,
		},
		License: LicenseConfig{
			LicenseFile: "data/license.key",
		},
		Security: SecurityConfig{
			JWTExpireDays:    7,
			MaxLoginAttempts: 5,
			LockoutMinutes:   5,
		},
		Backup: BackupConfig{
			Enabled:       true,
			IntervalHours: 24,
			KeepCount:     7,
			BackupPath:    getDefaultBackupPath(),
		},
	}
}

// Load reads config from path, applying defaults for missing fields
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c := Default()
			c.SetPath(path)
			return c, nil
		}
		return nil, err
	}
	c := Default()
	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}
	c.SetPath(path)
	c.applyDefaults()
	return c, nil
}

// applyDefaults fills zero values
func (c *Config) applyDefaults() {
	d := Default()
	if c.Server.Port == 0 {
		c.Server.Port = d.Server.Port
	}
	if c.Server.Host == "" {
		c.Server.Host = d.Server.Host
	}
	if c.Server.AddressType == "" {
		c.Server.AddressType = d.Server.AddressType
	}
	if c.Server.DataDir == "" {
		c.Server.DataDir = d.Server.DataDir
	}
	if c.Server.DBPath == "" {
		c.Server.DBPath = d.Server.DBPath
	}
	if c.Server.MaxUploadMB == 0 {
		c.Server.MaxUploadMB = d.Server.MaxUploadMB
	}
	if c.Security.JWTExpireDays == 0 {
		c.Security.JWTExpireDays = d.Security.JWTExpireDays
	}
	if c.Security.MaxLoginAttempts == 0 {
		c.Security.MaxLoginAttempts = d.Security.MaxLoginAttempts
	}
	if c.Security.LockoutMinutes == 0 {
		c.Security.LockoutMinutes = d.Security.LockoutMinutes
	}
	if c.License.LicenseFile == "" {
		c.License.LicenseFile = d.License.LicenseFile
	}
	if c.Backup.IntervalHours == 0 {
		c.Backup.IntervalHours = d.Backup.IntervalHours
	}
	if c.Backup.KeepCount == 0 {
		c.Backup.KeepCount = d.Backup.KeepCount
	}
	if c.Backup.BackupPath == "" {
		c.Backup.BackupPath = d.Backup.BackupPath
	}
}

// Get returns the current config instance
func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return instance
}

var (
	mu       sync.RWMutex
	instance *Config
)

// Update applies new server settings and saves
func Update(fn func(*Config)) error {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = Default()
	}
	fn(instance)
	return nil
}

// SetInstance installs the active config
func SetInstance(c *Config) {
	mu.Lock()
	defer mu.Unlock()
	instance = c
}
