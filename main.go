package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mike-neutron/go_link_shortener/src/controllers"
	"github.com/mike-neutron/go_link_shortener/src/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".env")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

// @title Rates API
// @version 1.0
// @description api for getting rates, exchange rates
// @BasePath /
func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", controllers.HelloWorld)
	app.Post("/make", controllers.Make)
	app.Get("/:link", controllers.Get)
	log.Fatal(app.Listen(":8080"))
}
