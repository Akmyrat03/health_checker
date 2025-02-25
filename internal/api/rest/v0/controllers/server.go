package controllers

import (
	"checker/internal/api"
	"checker/internal/api/rest/v0/requests"
	"checker/internal/api/rest/v0/responses"
	app_errors "checker/internal/domain/app/errors"
	"checker/internal/domain/app/inputs"
	"checker/internal/domain/entities"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateServer(c *fiber.Ctx) error {
	var server requests.CreateServer

	if err := c.BodyParser(&server); err != nil {
		fmt.Printf("ERROR: Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(entities.Error{
			Loc:  []string{"body"},
			Msg:  "Failed to parse request body",
			Type: "validation_error",
		})
	}

	input := inputs.CreateServer{
		Name: server.Name,
		Url:  server.Url,
	}

	serverUseCase, err := api.MakeServerUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	id, err := serverUseCase.Create(c.Context(), input)
	if err != nil {
		fmt.Printf("ERROR: Failed to create server: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"createUseCase", "create"},
			Msg:  "Failed to create server",
			Type: "database_error",
		})
	}

	response := responses.CreateServer{
		Id: id,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func DeleteServer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		fmt.Printf("ERROR: Invalid id format: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(entities.Error{
			Loc:  []string{"form"},
			Msg:  "ServerID must be a valid id",
			Type: "bad_request",
		})
	}

	serverUseCase, err := api.MakeServerUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	err = serverUseCase.Delete(c.Context(), id)
	if err != nil {
		fmt.Printf("ERROR: Failed to delete Server: %v\n", err)
		if errors.Is(err, app_errors.ServerDoesNotExist) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(entities.Error{
				Loc:  []string{"serverUseCase", "delete"},
				Msg:  app_errors.ServerDoesNotExist.Error(),
				Type: "bad_request",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"serverUseCase", "delete"},
			Msg:  err.Error(),
			Type: "internal_server_error",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetServers(c *fiber.Ctx) error {
	serverUseCase, err := api.MakeServerUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	servers, err := serverUseCase.List(c.Context())
	if err != nil {
		fmt.Printf("ERROR: Failed to list servers: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"serverUseCase", "list"},
			Msg:  "Failed to list servers",
			Type: "database_error",
		})
	}

	response := servers

	return c.Status(fiber.StatusOK).JSON(response)
}
