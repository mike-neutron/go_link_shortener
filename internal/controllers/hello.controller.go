package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func HelloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋! It`s a link shortener in Golang.")
}
