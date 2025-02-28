package controllers

import (
	"checker/internal/api"
	"checker/internal/api/rest/v0/requests"
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GetBasicConfig godoc
// @Summary Get all basic configs
// @Desciption Get all necessary configs
// @Tags basic_config
// @Produce json
// @Success 200 {object} interface{} "success"
// @Failure 500 {object} entities.Error
// @Router /api/v0/basic [get]
func GetBasicConfig(c *fiber.Ctx) error {
	basicUseCase, err := api.MakeBasicUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	basic, err := basicUseCase.Get(c.Context())
	if err != nil {
		fmt.Printf("ERROR: Failed to get basic config datas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"basicUseCase", "get"},
			Msg:  "Failed to get basic config datas",
			Type: "database error",
		})
	}

	response := &basic

	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdateBasicConfig godoc
// @Summary Update basic config
// @Description Update necessary configs: check_interval, timeout and notification interval
// @Tags basic_config
// @Produce json
// @Param basic body requests.UpdateBasic true "Basic Config"
// @Success 200 {object} string
// @Router /api/v0/basic [put]
func UpdateBasicConfig(c *fiber.Ctx) error {
	var basic requests.UpdateBasic

	if err := c.BodyParser(&basic); err != nil {
		fmt.Printf("ERROR: Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(entities.Error{
			Loc:  []string{"body"},
			Msg:  "Failed to parse request body",
			Type: "validation_error",
		})
	}

	input := inputs.UpdateBasic{
		CheckInterval:        basic.CheckInterval,
		Timeout:              basic.Timeout,
		NotificationInterval: basic.NotificationInterval,
	}

	basicUseCase, err := api.MakeBasicUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	err = basicUseCase.Update(c.Context(), input)
	if err != nil {
		fmt.Printf("ERROR: Failed to update necessasy configs: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"basicUseCase", "update"},
			Msg:  err.Error(),
			Type: "internal_server_error",
		})
	}

	// responses := responses.BasicConfig{}

	return c.Status(fiber.StatusOK).JSON("ok")
}
