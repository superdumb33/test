package controllers

import (
	"rest-service/internal/dto"
	"rest-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(app *fiber.App, service *services.UserService) {
	handler := &UserController{service: service}

	app.Post("/api/v1/users/authorize", handler.Authorize)
	app.Post("/api/v1/users/refresh", handler.Refresh)
}

func (uc *UserController) Authorize(c *fiber.Ctx) error {
	var request dto.AuthorizeUserRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}
	if request.UUID.String() == "" {
		return c.SendStatus(400)
	}

	tokens, err := uc.service.Authorize(c.Context(), request.UUID, c.IP())
	if err != nil {
		return err
	}

	resp := dto.AuthorizeUserResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	return c.Status(200).JSON(resp)
}

func (uc *UserController) Refresh(c *fiber.Ctx) error {
	var request dto.RefreshTokensRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}
	if request.AccessToken == "" || request.RefreshToken == "" {
		return c.SendStatus(400)
	}
	
	tokens, err := uc.service.Refresh(c.Context(), request.AccessToken, request.RefreshToken, c.IP())
	if err != nil {
		return err
	}

	resp := dto.RefreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	return c.Status(200).JSON(resp)
}
