package service

import "github.com/gin-gonic/gin"

type ESService struct {
}

func (puller *ESService) Run() error {
	esResource := &ESResource{}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.Static("/assert", "./assert")
	r.GET("/index", esResource.index)
	r.GET("/search", esResource.search)
	r.GET("/users/:accountID/logs", esResource.logs)
	r.Run(":8080")
	return nil
}
