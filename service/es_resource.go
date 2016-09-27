package service

import (
	"crypto/rand"
	"fmt"
	"os"

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
	buffer, hits, err := pull.GenerateFile(result)

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
	buffer, hits, err := pull.GenerateFile(result)
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
