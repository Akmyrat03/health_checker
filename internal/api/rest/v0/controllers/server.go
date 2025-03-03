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

type HealthStatus struct {
	ServerName string `json:"server_name"`
	ServerURL  string `json:"server_url"`
	Status     string `json:"status"`
}

// @CreateServer godoc
// @Summary Create Server
// @Description Create Server
// @Tags servers
// @Produce json
// @Param name body requests.CreateServer true "Server Name"
// @Success 200 {object} responses.CreateServer "success"
// @Router /api/v0/servers [post]
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

// DeleteServer godoc
// @Summary Delete Server
// @Description Delete Server
// @Tags servers
// @Produce json
// @Param id query string true "Server ID"
// @Success 204 "Success"
// @Router /api/v0/servers [delete]
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

// GetServers godoc
// @Summary Get All Servers
// @Description Retrieve a list of all servers
// @Tags servers
// @Produce json
// @Success 200 {array} responses.CreateServer
// @Failure 500 {object} entities.Error
// @Router /api/v0/servers [get]
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

// func ShowStatus(c *fiber.Ctx) error {
// 	// Connect to the database
// 	pool, err := pgx.PostgresPool()
// 	if err != nil {
// 		log.Printf("[ShowStatus]: Failed to connect to database: %v\n", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"code":    fiber.StatusInternalServerError,
// 			"message": "Failed to connect to the database",
// 		})
// 	}

// 	// Initialize repositories
// 	serverRepo := pgx_repositories.NewPgxServerRepository(pool)
// 	basicRepo := pgx_repositories.NewPgxBasicRepository(pool)

// 	// Get servers from the database
// 	servers, err := serverRepo.List(c.Context())
// 	if err != nil {
// 		log.Printf("[ShowStatus]: Failed to fetch servers: %v\n", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"code":    fiber.StatusInternalServerError,
// 			"message": "Failed to fetch servers",
// 		})
// 	}

// 	// Get basic configuration from the database
// 	basicConfig, err := basicRepo.Get(c.Context())
// 	if err != nil {
// 		log.Printf("[ShowStatus]: Failed to fetch basic config: %v\n", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"code":    fiber.StatusInternalServerError,
// 			"message": "Failed to fetch basic config",
// 		})
// 	}

// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	status := []HealthStatus{}

// 	// Check each server concurrently
// 	for _, server := range servers {
// 		wg.Add(1)
// 		go func(server entities.Server) {
// 			defer wg.Done()

// 			err := usecases.CheckServer(c.Context(), server, basicConfig.Timeout)

// 			serverStatus := "Healthy"
// 			if err != nil {
// 				serverStatus = "Unhealthy"
// 			}

// 			mu.Lock()
// 			status = append(status, HealthStatus{
// 				ServerName: server.Name,
// 				ServerURL:  server.Url,
// 				Status:     serverStatus,
// 			})
// 			mu.Unlock()
// 		}(server)
// 	}

// 	wg.Wait()

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"servers": status,
// 	})
// }
