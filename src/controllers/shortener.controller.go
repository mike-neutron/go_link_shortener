package controllers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mike-neutron/go_link_shortener/src/initializers"
	"github.com/mike-neutron/go_link_shortener/src/models"
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

	var row models.Link
	err := initializers.DB.Where("original = ?", body.Link).First(&row)
	if err != nil {
		if err.RowsAffected != 1 {
			row = models.Link{
				Original: body.Link,
			}
			initializers.DB.Create(&row)
		}
	}

	id, _ := s.Encode([]uint64{uint64(row.ID)})
	return c.SendString(id)
}

func Get(c *fiber.Ctx) error {

	short := c.Params("link")

	ids := s.Decode(short)

	if len(ids) == 1 {
		id := ids[0]
		var row models.Link
		err := initializers.DB.Where("id = ?", id).First(&row).Error
		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(fmt.Sprint(row.Original))
	}

	return c.SendStatus(404)
}
