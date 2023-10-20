package controllers

import (
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mike-neutron/go_link_shortener/internal/initializers"
	"github.com/mike-neutron/go_link_shortener/internal/models"
	"github.com/sqids/sqids-go"
)

type MakeShortLinkRequest struct {
	Link string `json:"link" validate:"required,min=2,max=256"`
}

var s, _ = sqids.New()

func Make(c *fiber.Ctx) error {
	var (
		body MakeShortLinkRequest
		row  models.Link
	)
	if err := c.BodyParser(&body); err != nil {
		return c.SendStatus(400)
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.SendStatus(400)
	}

	url, err := url.Parse(body.Link)
	if err != nil {
		return c.SendStatus(400)
	}
	if len(url.Host) == 0 {
		url, err := url.Parse("//" + body.Link)
		if err != nil || len(url.Host) == 0 {
			return c.SendStatus(400)
		}
	}

	urlQuery := url.RawQuery
	urlPath := url.Path
	if len(url.RawQuery) > 0 {
		urlQuery = "?" + urlQuery
	}

	formattedLink := "//" + url.Host + urlPath + urlQuery
	query := initializers.DB.Where("original = ?", formattedLink).First(&row)
	if query.Error != nil {
		return c.SendStatus(500)
	}
	if query.RowsAffected != 1 {
		row = models.Link{
			Original: formattedLink,
		}
		initializers.DB.Create(&row)
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

		return c.SendString(row.Original)
	}

	return c.SendStatus(404)
}
