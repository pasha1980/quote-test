package service

import (
	"context"
	"github.com/google/uuid"
	"quote-app/internal"
)

type QuoteServiceImpl struct {
	repository internal.QuoteRepository
}

func (s *QuoteServiceImpl) GetSame(c context.Context, id string) ([]internal.Quote, error) {
	quote, err := s.repository.FindById(c, id)
	if err != nil {
		return nil, err
	}

	if quote == nil {
		return nil, internal.ErrNotFound
	}

	return s.repository.ListSameTo(c, quote)
}

func (s *QuoteServiceImpl) Like(c context.Context, id string) (*internal.Quote, error) {
	quote, err := s.repository.FindById(c, id)
	if err != nil {
		return nil, err
	}

	quote.Likes++
	_, err = s.repository.Update(c, quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (s *QuoteServiceImpl) GetRandom(c context.Context) (*internal.Quote, error) {
	return s.repository.FindRandom(c)
}

func (s *QuoteServiceImpl) Create(c context.Context, author string, content string) (*internal.Quote, error) {

	if exists, _ := s.repository.Exists(c, author, content); exists {
		return nil, internal.ErrAlreadyExists
	}

	return s.repository.Create(c, &internal.Quote{
		Id:      uuid.New().String(),
		Author:  author,
		Content: content,
	})
}

func NewQuoteService(repository internal.QuoteRepository) (*QuoteServiceImpl, error) {
	return &QuoteServiceImpl{
		repository: repository,
	}, nil
}
