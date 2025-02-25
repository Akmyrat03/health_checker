package controllers

import (
	"checker/internal/api/rest/v0/controllers"

	"github.com/gofiber/fiber/v2"
)

func GroupControllers(app *fiber.Router) {
	v0 := (*app).Group("/v0")

	protected := v0.Group("")
	// protected.Use(middleware.AuthMiddleware())

	protected.Post("/servers", controllers.CreateServer)
	protected.Delete("/servers", controllers.DeleteServer)
	protected.Get("/servers", controllers.GetServers)

	protected.Get("/basic", controllers.GetBasicConfig)

	v0.Get("/health", controllers.ShowStatus)
}
