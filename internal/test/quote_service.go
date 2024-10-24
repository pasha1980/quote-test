package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"quote-app/infrastructure/elasticsearch"
	"quote-app/internal"
)

type QuoteServiceTest struct {
	suite.Suite
	service internal.QuoteService
}

func (t *QuoteServiceTest) SetupTest() {
	elasticsearch.Client.Indices.Delete([]string{"quotes"})
	elasticsearch.Client.Indices.Create("quotes")
	doc1 := internal.Quote{
		Id:      "5a64ed63-4a08-4f36-986a-440f22e60332",
		Likes:   1,
		Author:  "William Shakespeare",
		Content: "All the world is a stage, And all the men and women merely players. They have their exits and entrances; Each man in his time plays many parts.",
	}
	doc1Json, _ := json.Marshal(doc1)
	elasticsearch.Client.Create("quotes", "5a64ed63-4a08-4f36-986a-440f22e60332", bytes.NewReader(doc1Json))

	doc2 := internal.Quote{
		Id:      "1b83a7a8-60f0-4d9f-b610-e661fa3aafe6",
		Likes:   0,
		Author:  "Tony Robbins",
		Content: "We can change our lives. We can do, have, and be exactly what we wish.",
	}
	doc2Json, _ := json.Marshal(doc2)
	elasticsearch.Client.Create("quotes", "5a64ed63-4a08-4f36-986a-440f22e60332", bytes.NewReader(doc2Json))
}

func (t *QuoteServiceTest) TestCreateQuote() {
	quote, err := t.service.Create(context.Background(), "Margaret Thatcher", "To wear your heart on your sleeve isn't a very good plan; you should wear it inside, where it functions best.")
	if err != nil {
		t.Fail(err.Error())
	}

	t.Assert().NotNil(quote)
}

func (t *QuoteServiceTest) TestLikeQuote() {
	quote, err := t.service.Like(context.Background(), "5a64ed63-4a08-4f36-986a-440f22e60332")
	if err != nil {
		t.Fail(err.Error())
	}

	t.Assert().Same(int64(2), quote.Likes)
}

func (t *QuoteServiceTest) TestGetRandomQuote() {
	quote, err := t.service.GetRandom(context.Background())
	if err != nil {
		t.Fail(err.Error())
	}

	t.Assert().NotNil(quote)
}

func (t *QuoteServiceTest) TestGetSameQuote() {
	quote, err := t.service.GetSame(context.Background(), "5a64ed63-4a08-4f36-986a-440f22e60332")
	if err != nil {
		t.Fail(err.Error())
	}

	t.Assert().NotNil(quote)
}

func NewQuoteServiceTest(service internal.QuoteService) *QuoteServiceTest {
	return &QuoteServiceTest{
		service: service,
	}
}
