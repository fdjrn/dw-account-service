package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/dw-account-service/internal/db/entity"
	"github.com/dw-account-service/internal/utilities"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendDefaultErrResponse(prefix string, err error, c *fiber.Ctx) error {

	utilities.Log.Println(err)

	errMsg := err.Error()
	if errors.Is(err, context.DeadlineExceeded) {
		errMsg = "context deadline exceeded"
	}

	if errors.Is(err, mongo.ErrClientDisconnected) {
		errMsg = "client disconnected"
	}

	if mongo.IsNetworkError(err) {
		errMsg = "error connection occurred"
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		errMsg = "entity not found"
	}

	return c.Status(500).JSON(entity.Responses{
		Success: false,
		Message: fmt.Sprintf("%s%s", prefix, errMsg),
		Data:    nil,
	})
}

func SendDefaultPaginationErrResponse(prefix string, err error, c *fiber.Ctx) error {

	errMsg := err.Error()
	if errors.Is(err, context.DeadlineExceeded) {
		errMsg = "context deadline exceeded"
	}

	return c.Status(500).JSON(entity.PaginatedResponse{
		Success: false,
		Message: fmt.Sprintf("%s%s", prefix, errMsg),
		Data:    entity.PaginatedDetailResponse{},
	})
}
