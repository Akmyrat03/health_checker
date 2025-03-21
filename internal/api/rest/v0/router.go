package controllers

import (
	"checker/internal/api/rest/v0/controllers"

	"github.com/gofiber/fiber/v2"
)

func GroupControllers(app *fiber.Router) {
	v0 := (*app).Group("/v0")

	protected := v0.Group("")

	protected.Post("/servers", controllers.CreateServer)
	protected.Delete("/servers", controllers.DeleteServer)
	protected.Get("/servers", controllers.GetServers)

	protected.Get("/basic", controllers.GetBasicConfig)
	protected.Put("/basic", controllers.UpdateBasicConfig)

	protected.Post("/receiver", controllers.CreateReceiver)
	protected.Delete("/receiver", controllers.DeleteReceiver)
	protected.Get("/receiver", controllers.GetReceivers)
	protected.Get("/receiver/mute", controllers.MuteReceiver)
	protected.Get("/receiver/unmute", controllers.UnmuteReceiver)
}
