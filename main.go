package main

import (
	"fmt"
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
	date_gte := c.Query("date_gte")
	date_lte := c.Query("date_lte")
	fmt.Println("###", date_gte, date_lte)
	pull, err := NewPuller(date_gte, date_lte)
	if err != nil {
		panic(err)
	}
	result, err := pull.Search()
	file, err := pull.GenerateFile(result)
	fmt.Println("############", file.Name)
	c.File("/tmp/hello")
}
