package http

import (
	"gopastebin/config"
	"gopastebin/internal/controller"
	"gopastebin/internal/http/middleware"
	"gopastebin/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	authSrv *service.AuthService,
	pasteSrv *service.PasteService,
	userSrv *service.UserService,
	settingsController *SettingsController,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(60))

	authCtrl := controller.NewAuthController(authSrv)
	pasteCtrl := controller.NewPasteController(pasteSrv)
	userCtrl := controller.NewUserController(userSrv)

	api := r.Group("/api/v1")

	// health
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// public
	api.POST("/auth/register", authCtrl.Register)
	api.POST("/auth/login", authCtrl.Login)
	api.GET("/auth/me", authCtrl.Me)
	api.GET("/auth/verify", authCtrl.Verify)

	api.GET("/pastes/:slug", pasteCtrl.Get)
	api.GET("/pastes/:slug/raw", pasteCtrl.Raw)
	api.GET("/pastes/recent", pasteCtrl.Recent)
	api.POST("/pastes/anon", pasteCtrl.CreateAnon)

	api.GET("/users/:username/public", userCtrl.PublicProfile)

	// private
	priv := api.Group("/")
	priv.Use(middleware.Auth(authSrv))
	{
		priv.POST("/pastes", pasteCtrl.Create)
		priv.PUT("/pastes/:slug", pasteCtrl.Update)
		priv.DELETE("/pastes/:slug", pasteCtrl.Delete)

		priv.GET("/me/pastes", pasteCtrl.MyPastes)
		priv.GET("/me/folders", pasteCtrl.MyFolders)
		priv.GET("/me/pastes/by-folder", pasteCtrl.MyPastesByFolder)

		priv.GET("/users/me", userCtrl.Me)
		priv.PUT("/users/me", userCtrl.UpdateMe)

		priv.GET("/settings/storage", settingsController.GetStorage)
		priv.POST("/settings/storage", middleware.Admin(), settingsController.SetStorage)
	}

	return r
}
