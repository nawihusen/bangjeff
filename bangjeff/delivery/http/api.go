package http

import (
	"bangjeff/bangjeff/delivery/http/handler"
	"bangjeff/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// RouterAPI is main router for this Service User REST API
func RouterAPI(app *fiber.App, user domain.UserUsecase) {
	userHandler := &handler.UserHandler{UserUsecase: user}

	basePath := viper.GetString("server.base_path")
	path := app.Group(basePath)

	// Article Management
	path.Post("/user/signin", userHandler.SignIn)
	path.Post("/user/signup", userHandler.SignUp)
	path.Get("/user", userHandler.Authorization(), userHandler.GetUsers)
	path.Get("/user/profile", userHandler.Authorization(), userHandler.GetProfile)
	path.Get("/user/:id", userHandler.Authorization(), userHandler.GetUser)
}
