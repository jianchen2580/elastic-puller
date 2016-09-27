package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	//"net/http"
)

func main() {
	r := gin.Default()
	r.SetHTMLTemplate(html)
	r.Static("/static", "./static")
	r.GET("/index", index)
	r.GET("/search", search)
	r.GET("/users/:accountID/logs", logs)
	r.Run(":8080")
}

func index(c *gin.Context) {
	c.HTML(200, "ES_Puller", gin.H{})
}

func logs(c *gin.Context) {
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
	buffer, err := pull.GenerateFile(result)
	c.String(200, buffer.String())
}

func search(c *gin.Context) {
	index := c.Query("index")
	dateGte := c.Query("date_gte")
	dateLte := c.Query("date_lte")
	accountID := c.Query("account_id")
	appID := c.Query("app_id")
	sessionID := c.Query("session_id")
	//fmt.Println("###", date_gte, date_lte)
	pull, err := NewPuller(index, dateGte, dateLte, accountID, sessionID, appID)
	if err != nil {
		panic(err)
	}
	result, err := pull.Search()
	buffer, err := pull.GenerateFile(result)
	c.String(200, buffer.String())
}
