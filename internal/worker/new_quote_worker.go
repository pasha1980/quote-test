package worker

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"quote-app/internal"
	"time"
)

type AddQuoteWorker struct {
	service internal.QuoteService
}

func (w *AddQuoteWorker) run() error {
	resp, err := http.Get("https://api.quotable.io/random")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var quoteResp struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	respBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &quoteResp)
	if err != nil {
		return err
	}

	_, err = w.service.Create(context.Background(), quoteResp.Author, quoteResp.Content)
	if err != nil {
		return err
	}
	return nil
}

func InitAddQuoteWorker(service internal.QuoteService) {
	go func() {
		var worker = &AddQuoteWorker{
			service: service,
		}
		for range time.Tick(time.Second) {
			err := worker.run()
			if err != nil {
				log.Println(err)
			}
		}
	}()
}
