package main

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"os"
	"quote-app/infrastructure/elasticsearch"
	"quote-app/infrastructure/logger"
	"quote-app/internal"
	"quote-app/internal/repository"
	"quote-app/internal/service"
	"quote-app/internal/test"
	"testing"
)

func TestApp(t *testing.T) {
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
				test.NewQuoteServiceTest,
			),
			fx.Invoke(
				elasticsearch.Init,
				func(serviceTest *test.QuoteServiceTest) {
					suite.Run(t, serviceTest)
				},
			),
		),
	)
}
