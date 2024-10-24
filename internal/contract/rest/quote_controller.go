package rest

import (
	"github.com/gofiber/fiber/v2"
	"quote-app/internal"
)

type QuoteController struct {
	service internal.QuoteService
}

func (c *QuoteController) GetRandom(ctx *fiber.Ctx) error {
	quote, err := c.service.GetRandom(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(quote)
}

func (c *QuoteController) Like(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	quote, err := c.service.Like(ctx.Context(), idParam)
	if err != nil {
		return err
	}

	return ctx.JSON(quote)
}

func (c *QuoteController) Same(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	quotes, err := c.service.GetSame(ctx.Context(), idParam)
	if err != nil {
		return err
	}

	return ctx.JSON(quotes)
}

func (c *QuoteController) Routes(app *fiber.App) {
	app.Get("/quote/random", c.GetRandom)
	app.Post("/quote/:id/like", c.Like)
	app.Get("/quote/:id/same", c.Same)
}

func NewQuoteController(service internal.QuoteService) (*QuoteController, error) {

	return &QuoteController{
		service: service,
	}, nil
}
