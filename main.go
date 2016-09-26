package main

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	//"net/http"
)

func main() {
	r := gin.Default()
	r.SetHTMLTemplate(html)

	r.GET("/index", index)
	r.GET("/search", search)
	r.Run(":8080")
}

func index(c *gin.Context) {
	c.HTML(200, "ES_Puller", gin.H{})
}

func search(c *gin.Context) {
	dateGte := c.Query("date_gte")
	dateLte := c.Query("date_lte")
	accountID := c.Query("account_id")
	appID := c.Query("app_id")
	sessionID := c.Query("session_id")
	//fmt.Println("###", date_gte, date_lte)
	pull, err := NewPuller(dateGte, dateLte, accountID, sessionID, appID)
	if err != nil {
		panic(err)
	}
	result, err := pull.Search()
	buffer, err := pull.GenerateFile(result)
	c.String(200, buffer.String())
}
