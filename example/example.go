package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blevesearch/bleve"
	_ "github.com/sunvim/bleve-sego"
)

func main() {
	Example()
}

func Example() {
	INDEX_DIR := "sego.bleve"
	messages := []struct {
		Id   string
		Body string
	}{
		{
			Id:   "1",
			Body: "你好",
		},
		{
			Id:   "2",
			Body: "交代",
		},
		{
			Id:   "3",
			Body: "闵浦二桥",
		},
	}

	indexMapping := bleve.NewIndexMapping()
	os.RemoveAll(INDEX_DIR)
	// clean index when example finished
	defer os.RemoveAll(INDEX_DIR)

	err := indexMapping.AddCustomTokenizer("sego",
		map[string]interface{}{
			"dictpath": "d:/gowork/src/github.com/huichen/sego/data/dictionary.txt",
			"type":     "sego",
		},
	)
	if err != nil {
		panic(err)
	}
	err = indexMapping.AddCustomAnalyzer("sego",
		map[string]interface{}{
			"type":      "sego",
			"tokenizer": "sego",
		},
	)
	if err != nil {
		panic(err)
	}
	indexMapping.DefaultAnalyzer = "sego"

	index, err := bleve.New(INDEX_DIR, indexMapping)
	if err != nil {
		panic(err)
	}
	for _, msg := range messages {
		if err := index.Index(msg.Id, msg); err != nil {
			panic(err)
		}
	}

	querys := []string{
		"你好世界",
		"亲口交代",
		"二桥",
	}

	for _, q := range querys {
		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		req.Highlight = bleve.NewHighlight()
		res, err := index.Search(req)
		if err != nil {
			panic(err)
		}
		fmt.Println(prettify(res))
	}

	// Output:
	// [{"id":"1","score":0.27650412875470115}]
	// [{"id":"2","score":0.27650412875470115}]
	// [{"id":"3","score":0.8888941720598743}]
}

func prettify(res *bleve.SearchResult) string {
	type Result struct {
		Id    string  `json:"id"`
		Score float64 `json:"score"`
	}
	results := []Result{}
	for _, item := range res.Hits {
		results = append(results, Result{item.ID, item.Score})
	}
	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return string(b)
}
