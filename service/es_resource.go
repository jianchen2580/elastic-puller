package service

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v3"
	//"net/http"
	//"io/ioutil"
)

type ESResource struct {
}

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
	fmt.Println("###########", p.TimeStart, p.TimeEnd, p.SessionID)
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
	// TODO: hardcode index name here, replace with real later
	// TODO: the size of return record should be limited
	searchResult, err := client.Search().
		Index("log").        // search in index "log"
		Query(query).        // specify the query
		Sort("date", true).  // sort by "user" field, ascending
		From(0).Size(10000). // take documents 0-10000
		Pretty(true).        // pretty print request and response JSON
		Do()                 // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	// fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// TotalHits is another convenience function that works even when something goes wrong.
	// fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())
	return searchResult, nil
}

func (p *Puller) GenerateResult(sr *elastic.SearchResult) (bytes.Buffer, int64, error) {
	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Log
	var buffer bytes.Buffer
	for _, item := range sr.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Log); ok {
			s := fmt.Sprintf("%s %s %s %s\n", t.Timestamp, t.Program, t.Host, t.Message)
			buffer.WriteString(s)
		}
	}
	return buffer, sr.TotalHits(), nil
}

func (er *ESResource) index(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{})
}

func (er *ESResource) logs(c *gin.Context) {
	accountID := c.Param("accountID")
	//index := c.Query("index")
	index := "log"
	dateGte := c.Query("date_gte")
	dateLte := c.Query("date_lte")
	appID := c.Query("app_id")
	sessionID := c.Query("session_id")
	pull, err := NewPuller(index, dateGte, dateLte, accountID, sessionID, appID)
	if err != nil {
		panic(err)
	}
	result, err := pull.Search()
	buffer, hits, err := pull.GenerateResult(result)

	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		panic(err)
	}
	uuid := fmt.Sprintf("%X%X%X%X", b[0:4], b[4:8], b[8:12], b[10:])
	logFile := uuid + ".log"
	if err != nil {
		panic(err)
	}
	fo, err := os.Create("./static/" + logFile)
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	buffer.WriteTo(fo)
	c.JSON(200, gin.H{
		"hits":    hits,
		"logfile": "http://localhost:8080/static/" + logFile,
	})
}

func (er *ESResource) search(c *gin.Context) {
	index := c.Query("index")
	dateGte := c.Query("date_gte")
	dateLte := c.Query("date_lte")
	accountID := c.Query("account_id")
	appID := c.Query("app_id")
	sessionID := c.Query("session_id")
	pull, err := NewPuller(index, dateGte, dateLte, accountID, sessionID, appID)
	if err != nil {
		panic(err)
	}
	result, err := pull.Search()
	buffer, hits, err := pull.GenerateResult(result)
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		panic(err)
	}
	uuid := fmt.Sprintf("%X%X%X%X", b[0:4], b[4:8], b[8:12], b[10:])
	logFile := uuid + ".log"
	if err != nil {
		panic(err)
	}
	fo, err := os.Create("./static/" + logFile)
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	buffer.WriteTo(fo)
	c.HTML(200, "logs.tmpl", gin.H{
		"hits":    hits,
		"logfile": logFile,
	})
}
