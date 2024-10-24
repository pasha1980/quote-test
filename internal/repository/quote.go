package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"quote-app/infrastructure/elasticsearch"
	"quote-app/internal"
	"strings"
)

const quotesIndex = "quotes"

type QuoteRepositoryImpl struct {
}

func (r *QuoteRepositoryImpl) Create(c context.Context, quote *internal.Quote) (*internal.Quote, error) {
	docJson, err := json.Marshal(quote)
	if err != nil {
		return nil, err
	}

	resp, err := elasticsearch.Client.Create(quotesIndex, quote.Id, bytes.NewReader(docJson))
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	return quote, nil
}

func (r *QuoteRepositoryImpl) Update(c context.Context, quote *internal.Quote) (*internal.Quote, error) {
	doc := elasticsearch.MapToDoc(quote)
	docJson, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	resp, err := elasticsearch.Client.Update(quotesIndex, quote.Id, bytes.NewReader(docJson))
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	return quote, nil
}

func (r *QuoteRepositoryImpl) FindRandom(c context.Context) (*internal.Quote, error) {
	query := `
{
  "size": 1,
  "sort": [
    {
      "_script": {
        "script": "Math.random() * 200000",
        "type": "number",
        "order": "asc"
      }
    },
    {
      "likes": "desc"
    }
  ]
}`
	resp, err := elasticsearch.Client.Search(
		elasticsearch.Client.Search.WithIndex(quotesIndex),
		elasticsearch.Client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	quotes, err := elasticsearch.MapSearch[internal.Quote](resp.Body)
	if len(quotes) == 0 {
		return nil, internal.ErrNotFound
	}

	return &quotes[0], err
}

func (r *QuoteRepositoryImpl) FindById(c context.Context, id string) (*internal.Quote, error) {
	resp, err := elasticsearch.Client.Get(quotesIndex, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	return elasticsearch.MapItem[internal.Quote](resp.Body)
}

func (r *QuoteRepositoryImpl) Exists(c context.Context, author string, content string) (bool, error) {
	query := map[string]any{
		"size": 1,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"term": map[string]string{
							"author.keyword": author,
						},
					},
					{
						"term": map[string]string{
							"content.keyword": content,
						},
					},
				},
			},
		},
	}
	queryJson, _ := json.Marshal(query)
	resp, err := elasticsearch.Client.Search(
		elasticsearch.Client.Search.WithIndex(quotesIndex),
		elasticsearch.Client.Search.WithBody(bytes.NewReader(queryJson)),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	quotes, err := elasticsearch.MapSearch[internal.Quote](resp.Body)
	if len(quotes) == 0 {
		return false, nil
	}

	return true, nil
}

func (r *QuoteRepositoryImpl) ListSameTo(c context.Context, quote *internal.Quote) ([]internal.Quote, error) {
	query := map[string]any{
		"sort": []any{
			"_score",
			map[string]string{
				"likes": "desc",
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"should": []map[string]any{
					{
						"match": map[string]any{
							"author": map[string]any{
								"query":     quote.Author,
								"boost":     2,
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"match": map[string]any{
							"content": map[string]any{
								"query":     quote.Content,
								"fuzziness": "AUTO",
							},
						},
					},
				},
				"must_not": []map[string]any{
					{
						"term": map[string]string{
							"id.keyword": quote.Id,
						},
					},
				},
			},
		},
	}
	queryJson, _ := json.Marshal(query)
	resp, err := elasticsearch.Client.Search(
		elasticsearch.Client.Search.WithIndex(quotesIndex),
		elasticsearch.Client.Search.WithBody(bytes.NewReader(queryJson)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return elasticsearch.MapSearch[internal.Quote](resp.Body)
}

func NewQuoteRepository() (*QuoteRepositoryImpl, error) {
	elasticsearch.Client.Indices.Create(quotesIndex)
	return &QuoteRepositoryImpl{}, nil
}
