package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/gofiber/swagger"
	_ "github.com/mike-neutron/go_link_shortener/docs"
	"github.com/mike-neutron/go_link_shortener/internal/controllers"
	"github.com/mike-neutron/go_link_shortener/internal/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".env")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

// @title Link shortener API
// @version 1.0
// @description api for making short links
// @BasePath /
func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", controllers.HelloWorld)
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/api/make", controllers.Make)
	app.Get("/api/get/:short<maxLen(100);alpha>", controllers.Get)
	log.Fatal(app.Listen(":8080"))
}
