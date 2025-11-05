package config

import "os"

type Config struct {
	AppPort          string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTRefreshSecret string

	DefaultStorage string

	// firebase
	FirebaseCredsPath string
	FirebaseBucket    string
	FirebaseFolder    string

	// where to store current storage ("firebase"/"local")
	StorageSettingsPath string

	// smtp
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string

	// fallback local
	LocalStoragePath string
}

func env(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func Load() *Config {
	return &Config{
		AppPort:          env("APP_PORT", "8080"),
		DBHost:           env("DB_HOST", "localhost"),
		DBPort:           env("DB_PORT", "5432"),
		DBUser:           env("DB_USER", "postgres"),
		DBPassword:       env("DB_PASSWORD", "postgres"),
		DBName:           env("DB_NAME", "gopaste"),
		JWTSecret:        env("JWT_SECRET", "secret"),
		JWTRefreshSecret: env("JWT_REFRESH_SECRET", "secret_refresh"),

		DefaultStorage: env("DEFAULT_STORAGE", "firebase"),

		// ТУТ уже правильный дефолт – сервисный ключ
		FirebaseCredsPath:   env("FIREBASE_CREDENTIALS", "/app/secrets/firebase-admin.json"),
		FirebaseBucket:      env("FIREBASE_BUCKET", ""),
		FirebaseFolder:      env("FIREBASE_FOLDER", "pastes"),
		StorageSettingsPath: env("STORAGE_SETTINGS_PATH", "/app/data/firebase_settings.json"),

		SMTPHost: env("SMTP_HOST", ""),
		SMTPPort: env("SMTP_PORT", "587"),
		SMTPUser: env("SMTP_USER", ""),
		SMTPPass: env("SMTP_PASSWORD", ""),

		LocalStoragePath: env("LOCAL_STORAGE_PATH", "/app/storage-data"),
	}
}
