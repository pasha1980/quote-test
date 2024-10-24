package elasticsearch

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
	"log"
	"quote-app/config"
)

type Doc struct {
	Doc any `json:"doc"`
}

type searchResponse struct {
	Hits struct {
		Hits []itemResponse `json:"hits"`
	} `json:"hits"`
}

type itemResponse struct {
	Source json.RawMessage `json:"_source"`
}

var Client *elasticsearch.Client

func Init() {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			config.Get().ElasticSearchURL,
		},
		APIKey: config.Get().ElasticSearchKey,
	})
	if err != nil {
		panic(err)
	}

	Client = client
}

func MapSearch[T any](response io.ReadCloser) ([]T, error) {
	body, _ := io.ReadAll(response)

	var searchResult searchResponse
	err := json.Unmarshal(body, &searchResult)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var res []T
	for _, rawItem := range searchResult.Hits.Hits {
		var entity T
		err = json.Unmarshal(rawItem.Source, &entity)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = append(res, entity)
	}

	return res, nil
}

func MapItem[T any](response io.ReadCloser) (*T, error) {
	body, _ := io.ReadAll(response)

	var item itemResponse
	err := json.Unmarshal(body, &item)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var entity T
	err = json.Unmarshal(item.Source, &entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &entity, nil
}

func MapToDoc(body any) *Doc {
	return &Doc{
		Doc: body,
	}
}
