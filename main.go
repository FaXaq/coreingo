package main

import (
	"github.com/FaXaq/gjp"
	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"log"
)

//global config variables
var (
	CallbackStart string
	CallbackEnd   string
	WorkPath      string
	LogPath       string
	Port          string
)

func main() {
	err := GetEnvConfig()

	if err != nil {
		log.Fatal(err.Error()) //if not config then break
	}

	//create job pool
	jobPool := gjp.New(3)
	jobPool.Start()

	api := iris.New()

	restConfig := config.Rest{
		IndentJSON: true,
	}

	api.Config().Render.Rest = restConfig

	api.Get("/", func(c *iris.Context) {
		TestPing(c)
	})

	api.Get("/jobs/progress", func(c *iris.Context) {
		GetMyJobProgress(c, jobPool)
	})

	api.Post("/jobs", func(c *iris.Context) {
		CreateJob(c, jobPool)
	})

	api.Get("/jobs/search", func(c *iris.Context) {
		SearchJob(c, jobPool)
	})

	api.Listen(":" + Port)
}
