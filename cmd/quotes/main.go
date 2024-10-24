package main

import (
	"go.uber.org/fx"
	"os"
	"quote-app/infrastructure/elasticsearch"
	"quote-app/infrastructure/fiber"
	"quote-app/infrastructure/logger"
	"quote-app/internal"
	"quote-app/internal/contract/rest"
	"quote-app/internal/repository"
	"quote-app/internal/service"
	"quote-app/internal/worker"
)

func main() {
	logger.Init(os.Stdout)
	fx.New(
		fx.Options(
			fx.Provide(
				fx.Annotate(
					service.NewQuoteService,
					fx.As(new(internal.QuoteService)),
				),
				fx.Annotate(
					repository.NewQuoteRepository,
					fx.As(new(internal.QuoteRepository)),
				),
				rest.NewQuoteController,
			),
			fx.Invoke(
				elasticsearch.Init,
				fiber.Init,
				worker.InitAddQuoteWorker,
			),
		),
	).Run()
}
