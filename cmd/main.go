package main

import (
	"log"
	"rest-service/internal/controllers"
	"rest-service/internal/infrastructure/database"
	"rest-service/internal/infrastructure/repository/pgxrepo"
	"rest-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	db := database.ConnectToDB()

	userRepo := pgxrepo.NewPgxUserRepository(db)
	userService := services.NewUserService(userRepo)
	controllers.NewUserController(app, userService)



	log.Fatal(app.Listen(":3000"))
}