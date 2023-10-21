package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mike-neutron/go_link_shortener/internal/initializers"
	"github.com/mike-neutron/go_link_shortener/internal/models"
	"github.com/mike-neutron/go_link_shortener/internal/services"
)

type MakeRequest struct {
	Short string `json:"short" example:"da3rsf" validate:"required,min=6,max=100"`
	Link  string `json:"original" example:"http://example.com/" validate:"required,min=1,max=1000"`
}

type MakeResponse struct {
	Short string `json:"short" example:"da3rsf"`
}

type GetResponse struct {
	Original string `json:"original" example:"http://example.com/"`
}

//	@Summary		Make short link
//	@Description	Make short link
//	@Produce		json
//  @Param request body controllers.MakeRequest true "Make request"
//	@Success		200	{object} MakeResponse
//  @Failure        400
//	@Router			/api/make [post]
func Make(c *fiber.Ctx) error {
	var (
		body MakeRequest
		row  models.Link
	)
	if err := c.BodyParser(&body); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	url, err := url.Parse(body.Link)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	if len(url.Host) == 0 {
		url, err := url.Parse("//" + body.Link)
		if err != nil || len(url.Host) == 0 {
			return c.SendStatus(http.StatusBadRequest)
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
		short = services.Shorten(uuid.New().ID())
	}

	row = models.Link{
		Short:    short,
		Original: formattedLink,
	}
	initializers.DB.Create(&row)

	return c.JSON(&MakeResponse{Short: short})
}

//	@Summary		Get original link by short equivalent
//	@Description	Get original link by short equivalent
//	@Produce		json
//	@Param			short	path		string	true "Short"	"String"
//	@Success		200	{object}	GetResponse
//  @Failure        400
//  @Failure        404
//	@Router			/api/get/{short} [get]
func Get(c *fiber.Ctx) error {

	short := c.Params("short")
	if len(short) == 0 {
		return c.SendStatus(http.StatusBadRequest)
	}

	var row models.Link
	err := initializers.DB.Where("short = ?", short).First(&row).Error
	fmt.Print(err)

	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	return c.JSON(&GetResponse{Original: row.Original})
}
