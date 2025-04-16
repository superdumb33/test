package main

import (
	"log"
	"os"
	"rest-service/internal/controllers"
	"rest-service/internal/infrastructure/database"
	"rest-service/internal/infrastructure/repository/pgxrepo"
	"rest-service/internal/infrastructure/smtpclient"
	"rest-service/internal/services"


	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	db := database.ConnectToDB()
	smtpClient := smtpclient.NewSMTPClient(os.Getenv("SMTP_FROM"), os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"), "", "")

	authRepo := pgxrepo.NewPgxAuthRepository(db)
	authService := services.NewUserService(authRepo, smtpClient)
	controllers.NewAuthController(app, authService)



	log.Fatal(app.Listen(":3000"))
}