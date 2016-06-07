package main

import (
	"fmt"
	"github.com/FaXaq/gjp"
	"html"
	"log"
	"net/http"
)

//global config variables
var (
	CallbackStart string
	CallbackEnd   string
	WorkPath      string
	LogPath       string
)

func main() {

	err := GetEnvConfig()

	if err != nil {
		log.Fatal(err.Error())
	}

	jobPool := gjp.New(2) //create new jobPool with 2queues

	jobPool.Start()

	//start the jobPool
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Starting jobs")
		jobPool.Start()
	})

	//stop the jobPool
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		jobPool.ShutdownWorkPool()
	})

	//NYI
	// http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
	// 	list := jobPool.ListWaitingJobs()
	// 	fmt.Fprintf(w, "%v", list)
	// })

	//Create job
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q["command"] != nil &&
			q["fromFile"] != nil &&
			q["toFile"] != nil {

			//creating custom job
			newJob, err := NewJob(
				gjp.GenerateJobUID(),
				q["command"][0],
				q["fromFile"][0],
				q["toFile"][0])

			//if error while creating custom job, then print it in answer
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			} else {
				//if no error, queue the new job
				j := jobPool.QueueJob(newJob.Id, newJob.Name, newJob, 0)

				//get infos from the job
				jobjson, jsonerr := j.GetJobInfos()

				if jsonerr != nil {
					customJson := CreateCustomJson([]string{
						"Error",
						"Id",
					}, []string{
						jsonerr.Error(),
						j.GetJobId(),
					})
					fmt.Fprintf(w, "%v", string(customJson))
				} else {
					fmt.Fprintf(w, "%v", string(jobjson))
				}
			}
		} else {
			customJson := CreateCustomJson([]string{
				"Error",
			}, []string{
				"Couldn't create the job",
			})

			fmt.Fprintf(w, "%v", string(customJson))
		}
	})

	//get job infos
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q["id"] != nil {
			searchParam := q["id"][0]
			j, err := jobPool.GetJobFromJobId(searchParam)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				jobson, _ := j.GetJobInfos()
				fmt.Fprintf(w, "%v", string(jobson))
			}
		} else {
			fmt.Fprintf(w, "No search request")
		}
	})

	//get progress
	http.HandleFunc("/progress", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q["id"] != nil {
			searchParam := q["id"][0]
			j, err := jobPool.GetJobFromJobId(searchParam)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				timing, err := j.GetProgress(j.GetJobId())
				if err != nil {
					fmt.Fprintf(w, "%v", err.Error())
				} else {
					fmt.Fprintf(w, "%v", timing)
				}
			}
		} else {
			fmt.Fprintf(w, "No search request")
		}
	})

	//launch server
	log.Println(http.ListenAndServe("localhost:6060", nil))
}
