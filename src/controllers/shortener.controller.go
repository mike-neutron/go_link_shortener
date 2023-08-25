package controllers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sqids/sqids-go"
)

type MakeShortLinkRequest struct {
	Link string `json:"link" validate:"required,min=8,max=32"`
}

var s, _ = sqids.New()

func Make(c *fiber.Ctx) error {

	var body MakeShortLinkRequest
	if err := c.BodyParser(&body); err != nil {
		return c.SendStatus(400)
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.SendStatus(400)
	}

	id, _ := s.Encode([]uint64{1})
	return c.SendString(id)
}

func Get(c *fiber.Ctx) error {

	short := c.Params("link")

	ids := s.Decode(short)

	if len(ids) == 1 {
		id := ids[0]
		return c.SendString(fmt.Sprint(id))
	}

	return c.SendString("Not found")
}
