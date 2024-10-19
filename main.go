package main

import (
	"gopastebin/src/controllers"
	"gopastebin/src/db"
	"gopastebin/src/fb"
	"gopastebin/src/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	dbInstance, err := db.GetConnection()
	if err != nil {
		panic(err)
	}
	db.SetDB(dbInstance)

	err = fb.InitFirebaseApp()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(middleware.Logging())

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/confirm/:token", controllers.ConfirmEmail)
	r.GET("/paste/:slug", controllers.GetPaste)

	authorized := r.Group("/")
	authorized.Use(middleware.Auth())
	{
		authorized.POST("/upload", controllers.CreatePaste)
		authorized.POST("/paste", controllers.CreatePaste)
		authorized.PUT("/paste/:slug", controllers.UpdatePaste)
		authorized.GET("/user/:username/profile", controllers.GetUserProfile)
		authorized.PUT("/user/:username/profile", controllers.UpdateUserProfile)
	}

	r.Run(":8080")
}
