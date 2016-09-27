package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	//"net/http"
)

func main() {
	r := gin.Default()
	//r.SetHTMLTemplate(html)
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.GET("/index", index)
	r.GET("/search", search)
	r.GET("/users/:accountID/logs", logs)
	r.Run(":8080")
}

func index(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{})
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
	buffer, hits, err := pull.GenerateFile(result)
	c.HTML(200, "logs.tmpl", gin.H{
		"hits":   hits,
		"buffer": buffer,
	})
}

func search(c *gin.Context) {
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
	c.HTML(200, "logs.tmpl", gin.H{
		"hits":   hits,
		"buffer": buffer,
	})
}
