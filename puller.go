package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v3"
)

type Log struct {
	Timestamp string `json:"@timestamp"`
	Program   string `json:"program"`
	Host      string `json:"host"`
	Message   string `json:"message"`
}

type Puller struct {
	ESClient  *elastic.Client
	TimeStart string
	TimeEnd   string
	Session   string
}

func NewPuller(start string, end string) (*Puller, error) {
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	p := &Puller{
		ESClient:  client,
		TimeStart: start,
		TimeEnd:   end,
	}
	return p, nil
}

func (p *Puller) Search() (*elastic.SearchResult, error) {
	client := p.ESClient
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewRangeQuery("@timestamp").
		//Gte("2016-09-21T19:32:58.016Z").
		//Lte("now/d"))
		Gte(p.TimeStart).
		Lte(p.TimeEnd))
	query = query.Must(elastic.NewTermQuery("port", "35937"))
	// TODO: Print DSL, could remove in the future
	src, err := query.Source()
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	searchResult, err := client.Search().
		Index("log").       // search in index "twitter"
		Query(query).       // specify the query
		Sort("date", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do()                // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	//var ttyp Log
	//for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
	//	if t, ok := item.(Log); ok {
	//		fmt.Printf("log by %s %s %s %s\n", t.Timestamp, t.Program, t.Host, t.Message)
	//	}
	//}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())
	return searchResult, nil
}

func (p *Puller) GenerateFile(sr *elastic.SearchResult) (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "example")
	fmt.Println("filename is:", tmpfile.Name)
	if err != nil {
		panic(err)
	}

	var ttyp Log
	for _, item := range sr.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Log); ok {
			s := fmt.Sprintf("%s %s %s %s\n", t.Timestamp, t.Program, t.Host, t.Message)
			if _, err := tmpfile.Write([]byte(s)); err != nil {
				panic(err)
			}
		}
	}
	if err := tmpfile.Close(); err != nil {
		panic(err)
	}
	return tmpfile, nil
}
