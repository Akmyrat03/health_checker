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

// CreateReceiver godoc
// @Summary Create a receiver
// @Description Create a receiver
// @Tags receivers
// @Produce json
// @Param receiver_email body requests.CreateReceiver true "Receiver Email"
// @Success 200 {object} responses.CreateReceiver "success"
// @Failure 500 {object} entities.Error
// @Router /api/v0/receiver [post]
func CreateReceiver(c *fiber.Ctx) error {
	var receiver requests.CreateReceiver

	if err := c.BodyParser(&receiver); err != nil {
		fmt.Printf("ERROR: Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(entities.Error{
			Loc:  []string{"body"},
			Msg:  "Failed to parse request body",
			Type: "validation_error",
		})
	}

	input := inputs.CreateReceiver{
		Email: receiver.Email,
	}

	receiverUseCase, err := api.MakeReceiverUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	id, err := receiverUseCase.Create(c.Context(), input)
	if err != nil {
		fmt.Printf("ERROR: Failed to create receiver: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"receiverUseCase", "create"},
			Msg:  err.Error(),
			Type: "database_error",
		})
	}

	response := responses.CreateReceiver{
		ID: id,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteReceiver godoc
// @Summary Delete Receiver
// @Description Delete Receiver
// @Tags receivers
// @Produce json
// @Param id query string true "Receiver ID"
// @Success 204 "Success"
// @Failure 500 {object} entities.Error
// @Router /api/v0/receiver [delete]
func DeleteReceiver(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		fmt.Printf("ERROR: Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(entities.Error{
			Loc:  []string{"form"},
			Msg:  "ReceiverID must be a valid",
			Type: "bad_	request",
		})
	}

	receiverUseCase, err := api.MakeReceiverUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	err = receiverUseCase.Delete(c.Context(), id)
	if err != nil {
		fmt.Printf("ERROR: Failed to delete Receiver: %v\n", err)
		if errors.Is(err, app_errors.ReceiverDoesNotExist) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(entities.Error{
				Loc:  []string{"receiverUseCase", "delete"},
				Msg:  app_errors.ReceiverDoesNotExist.Error(),
				Type: "bad_request",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"receiverUseCase", "delete"},
			Msg:  err.Error(),
			Type: "internal_server_error",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)

}

// GetReceivers godoc
// @Summary Get all receivers
// @Description Retrieve a list of all receivers
// @Tags receivers
// @Produce json
// @Success 200 {object} responses.GetReceivers "success"
// @Failure 500 {object} entities.Error
// @Router /api/v0/receiver [get]
func GetReceivers(c *fiber.Ctx) error {
	receiverUseCase, err := api.MakeReceiverUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"server"},
			Msg:  err.Error(),
			Type: "processing_error",
		})
	}

	receivers, err := receiverUseCase.List(c.Context())
	if err != nil {
		fmt.Printf("ERROR: Failed to list receivers: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"receiverUseCase", "list"},
			Msg:  err.Error(),
			Type: "database_error",
		})
	}

	var response []responses.GetReceivers
	for _, muted := range receivers {
		response = append(response, responses.GetReceivers{
			ID:    muted.ID,
			Email: muted.Email,
			Muted: muted.Muted,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// MuteReceiver godoc
// @Summary Mute Receiver
// @Description Mute a receiver
// @Tags receivers
// @Produce json
// @Param email query string true "Receiver Email"
// @Success 200 {object} string "Receiver muted successfully"
// @Failure 400 {object} entities.Error "Bad Request"
// @Failure 500 {object} entities.Error "Internal server error"
// @Router /api/v0/receiver/mute [get]
func MuteReceiver(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.Error{
			Loc:  []string{"body", "email"},
			Msg:  "Email is required",
			Type: "validation_error",
		})
	}

	fmt.Println("MuteReceiver called with email:", email)

	receiverUseCase, err := api.MakeReceiverUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"server"},
			Msg:  err.Error(),
			Type: "processing_error",
		})
	}

	err = receiverUseCase.MuteStatus(c.Context(), email, true)
	if err != nil {
		fmt.Printf("ERROR: Failed to mute receiver: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"mute"},
			Msg:  err.Error(),
			Type: "database_error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Receiver muted successfully",
	})
}

// UnmuteReceiver godoc
// @Summary Unmute Receiver
// @Description Unmute a receiver
// @Tags receivers
// @Produce json
// @Param email query string true "Receiver Email"
// @Success 200 {object} string "Receiver unmuted successfully"
// @Failure 400 {object} entities.Error "Bad Request"
// @Failure 500 {object} entities.Error "Internal server error"
// @Router /api/v0/receiver/unmute [get]
func UnmuteReceiver(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(entities.Error{
			Loc:  []string{"body", "email"},
			Msg:  "Email is required",
			Type: "validation_error",
		})
	}

	receiverUseCase, err := api.MakeReceiverUseCase()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"server"},
			Msg:  err.Error(),
			Type: "processing_error",
		})
	}

	err = receiverUseCase.MuteStatus(c.Context(), email, false)
	if err != nil {
		fmt.Printf("ERROR: Failed to unmute receiver: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(entities.Error{
			Loc:  []string{"unmute"},
			Msg:  err.Error(),
			Type: "database_error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Receiver unmuted successfully",
	})
}
