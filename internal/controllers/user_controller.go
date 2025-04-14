package controllers

import (
	"rest-service/internal/dto"
	"rest-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service *services.UserService
}

func NewUserController (app *fiber.App, service *services.UserService) {
	handler := &UserController{service: service}

	app.Post("/api/v1/users/authorize", handler.Authorize)
}

func (uc *UserController) Authorize (c *fiber.Ctx) error {
	var request dto.AuthorizeUserRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}
	userIP := c.IP()

	tokens, err := uc.service.Authorize(request.UUID, userIP)
	if err != nil {
		return err
	}

	resp := dto.AuthorizeUserResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	return c.Status(200).JSON(resp)
}

