package internal

import "context"

type QuoteService interface {
	Like(c context.Context, id string) (*Quote, error)
	GetRandom(c context.Context) (*Quote, error)
	Create(c context.Context, author string, content string) (*Quote, error)
	GetSame(c context.Context, id string) ([]Quote, error)
}

type QuoteRepository interface {
	Create(c context.Context, quote *Quote) (*Quote, error)
	Update(c context.Context, quote *Quote) (*Quote, error)
	FindRandom(c context.Context) (*Quote, error)
	FindById(c context.Context, id string) (*Quote, error)
	Exists(c context.Context, author string, content string) (bool, error)
	ListSameTo(c context.Context, quote *Quote) ([]Quote, error)
}
