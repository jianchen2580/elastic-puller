package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"strconv"
	//"os"
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
	index     string
	TimeStart string
	TimeEnd   string
	SessionID string
	AccountID string
	AppID     string
}

func NewPuller(index string, start string, end string, accountID string, sessionID string, appID string) (*Puller, error) {
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	p := &Puller{
		ESClient:  client,
		TimeStart: start,
		TimeEnd:   end,
		SessionID: sessionID,
		AppID:     appID,
		AccountID: accountID,
	}
	return p, nil
}

func (p *Puller) Search() (*elastic.SearchResult, error) {
	client := p.ESClient
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewRangeQuery("@timestamp").
		Gte(p.TimeStart).
		Lte(p.TimeEnd))
	if len(p.AccountID) != 0 {
		query = query.Must(elastic.NewTermQuery("account_id", p.AccountID))
	}
	if len(p.AppID) != 0 {
		query = query.Must(elastic.NewTermQuery("app_id", p.AppID))
	}
	if len(p.SessionID) != 0 {
		query = query.Must(elastic.NewTermQuery("session_number", p.SessionID))
	}
	// TODO: Print DSL, could remove in the future
	src, err := query.Source()
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	searchResult, err := client.Search().
		Index("log").        // search in index "twitter"
		Query(query).        // specify the query
		Sort("date", true).  // sort by "user" field, ascending
		From(0).Size(10000). // take documents 0-9
		Pretty(true).        // pretty print request and response JSON
		Do()                 // execute
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

func (p *Puller) GenerateFile(sr *elastic.SearchResult) (bytes.Buffer, error) {
	tmpfile, err := ioutil.TempFile("", "example")
	fmt.Println("filename is:", tmpfile.Name)
	if err != nil {
		panic(err)
	}

	var ttyp Log
	var buffer bytes.Buffer
	for _, item := range sr.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Log); ok {
			s := fmt.Sprintf("%s %s %s %s\n", t.Timestamp, t.Program, t.Host, t.Message)
			buffer.WriteString(s)
		}
	}
	return buffer, nil
}
