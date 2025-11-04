package main

import (
	"context"
	"log"

	"gopastebin/config"
	"gopastebin/internal/domain"
	apphttp "gopastebin/internal/http"
	"gopastebin/internal/repository"
	"gopastebin/internal/service"
	"gopastebin/internal/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	// миграции доменных сущностей
	if err := repository.AutoMigrate(db); err != nil {
		log.Fatalf("migrate error: %v", err)
	}
	// миграция настроек
	if err := db.AutoMigrate(&domain.Setting{}); err != nil {
		log.Fatalf("settings migrate error: %v", err)
	}

	ctx := context.Background()

	// локальное хранилище есть всегда
	localFS := storage.NewLocalStore(cfg.LocalStoragePath)

	// пробуем поднять firebase
	var firebaseFS storage.FileStore
	firebaseFS, err = storage.NewFirebaseStore(ctx, cfg.FirebaseCredsPath, cfg.FirebaseBucket, cfg.FirebaseFolder)
	if err != nil {
		log.Printf("[warn] firebase init failed: %v — using only local", err)
		firebaseFS = nil
	}

	// репозитории
	userRepo := repository.NewUserRepository(db)
	pasteRepo := repository.NewPasteRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	// сервис настроек
	appSettings := service.NewAppSettingsService(settingRepo, cfg.DefaultStorage)

	// динамический стор
	smartFS := service.NewDynamicFileStore(appSettings, localFS, firebaseFS)

	// сервисы
	authService := service.NewAuthService(userRepo, cfg)
	pasteService := service.NewPasteService(pasteRepo, smartFS)
	userService := service.NewUserService(userRepo, pasteRepo)

	// контроллер настроек
	settingsController := apphttp.NewSettingsController(appSettings)

	// роутер
	r := apphttp.NewRouter(cfg, authService, pasteService, userService, settingsController)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
