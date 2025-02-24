package controllers

import (
	"checker/internal/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func GroupControllers(app *fiber.Router) {
	v0 := (*app).Group("/v0")
	{
		v0.Get("/health", controllers.ShowStatus)
	}

}
