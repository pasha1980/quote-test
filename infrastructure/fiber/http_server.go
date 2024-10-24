package fiber

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"quote-app/config"
	"quote-app/internal/contract/rest"
)

func Init(
	quoteController *rest.QuoteController,
) {

	app := fiber.New(fiber.Config{
		AppName:      config.Get().AppName,
		ErrorHandler: errorHandler,
	})

	app.Get("/", healthCheck)
	quoteController.Routes(app)

	go func(app *fiber.App) {
		app.Listen(config.Get().ServerAddress)
	}(app)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	log.Println(err)
	ctx.Status(http.StatusBadRequest)
	return ctx.JSON(map[string]string{
		"message": err.Error(),
	})
}

func healthCheck(ctx *fiber.Ctx) error {
	ctx.Status(http.StatusOK)
	return ctx.JSON(map[string]string{
		"message": "ok",
	})
}
