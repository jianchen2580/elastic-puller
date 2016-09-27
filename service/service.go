package service

import "github.com/gin-gonic/gin"

type ESService struct {
}

func (puller *ESService) Run() error {
	esResource := &ESResource{}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.GET("/index", index)
	r.GET("/search", search)
	r.GET("/users/:accountID/logs", logs)
	r.Run(":8080")
	return nil
}
