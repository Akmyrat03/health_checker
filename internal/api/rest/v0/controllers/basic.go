package controllers

import (
	"checker/internal/api"
	"checker/internal/domain/entities"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetBasicConfig(c *fiber.Ctx) error {
	basicUseCase, err := api.MakeBasicUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	basics, err := basicUseCase.List(c.Context())
	if err != nil {
		fmt.Printf("ERROR: Failed to get basic config datas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"basicUseCase", "get"},
			Msg:  "Failed to get basic config datas",
			Type: "database error",
		})
	}

	response := basics

	return c.Status(fiber.StatusOK).JSON(response)
}
