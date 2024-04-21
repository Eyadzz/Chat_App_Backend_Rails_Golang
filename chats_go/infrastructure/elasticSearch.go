package infrastructure

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
)

type ElasticSearch struct {
}

func NewElasticSearch() *ElasticSearch {
	return &ElasticSearch{}
}

func (ElasticSearch ElasticSearch) CreateIndex(name string) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	_, err = es.Indices.Create(name)
	if err != nil {
		log.Fatalf("Error creating the index: %s", err)
		return
	}

	log.Printf("Index created: %s", name)
}

func (ElasticSearch ElasticSearch) Insert(index string, message string) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    strings.NewReader(message),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()

	log.Printf("Document indexed: %s", res.String())
}

func (ElasticSearch ElasticSearch) Search(index string, query string) *esapi.Response {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	defer res.Body.Close()
	log.Printf("Search result: %s", res.String())
	return res

}
