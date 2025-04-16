package controllers

import (
	"rest-service/internal/dto"
	"rest-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(app *fiber.App, service *services.AuthService) {
	handler := &AuthController{service: service}

	app.Post("/api/v1/users/authorize", handler.Authorize)
	app.Post("/api/v1/users/refresh", handler.Refresh)
}

func (uc *AuthController) Authorize(c *fiber.Ctx) error {
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

func (uc *AuthController) Refresh(c *fiber.Ctx) error {
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
