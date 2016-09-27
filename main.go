package main

import (
	"log"
	"os"

	//"github.com/gin-gonic/gin"
	"github.com/jianchen2580/elastic-puller/service"

	"github.com/codegangsta/cli"
)

//"net/http"
//"io/ioutil"
func main() {
	app := cli.NewApp()
	app.Name = "es-puller"
	app.Usage = "work with `es-puller` microservice"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "config.yaml", Usage: "config file to use", EnvVar: "APP_CONFIG"},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				//cfg, err := getConfig(c)
				//if err != nil {
				//	log.Fatal(err)
				//	return
				//}
				svc := service.ESService{}
				if err := svc.Run(); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
