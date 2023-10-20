package controllers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mike-neutron/go_link_shortener/internal/initializers"
	"github.com/mike-neutron/go_link_shortener/internal/models"
	"github.com/sqids/sqids-go"
)

type MakeShortLinkRequest struct {
	Short string `json:"short"`
	Link  string `json:"link" validate:"required,min=2,max=1000"`
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
	if query.RowsAffected == 1 {
		return c.JSON(fiber.Map{"short": row.Short})
	}

	var short string
	if len(body.Short) > 0 {
		short = body.Short
	} else {
		short = Shorten(uuid.New().ID())
	}

	row = models.Link{
		Short:    short,
		Original: formattedLink,
	}
	initializers.DB.Create(&row)

	return c.JSON(fiber.Map{"short": row.Short})
}

func Get(c *fiber.Ctx) error {

	short := c.Params("short")
	if len(short) == 0 {
		return c.SendStatus(400)
	}

	var row models.Link
	err := initializers.DB.Where("short = ?", short).First(&row).Error
	fmt.Print(err)

	if err != nil {
		return c.SendStatus(400)
	}

	return c.JSON(fiber.Map{"link": row.Original})
}

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = uint32(len(alphabet))

func Shorten(id uint32) string {
	var (
		digits  []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}

	reverse(digits)

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
