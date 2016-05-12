package main

import (
	"fmt"
	"github.com/FaXaq/gjp"
	"os/exec"
	"net/http"
	"html"
	"log"
)

type (
	MyJob struct {
		name string
		commandName string
		arguments []string
	}
)

func (myjob *MyJob) ExecuteJob() (err *gjp.JobError) {
	out, commandError := exec.Command("echo","toto").Output()
	if commandError != nil {
		defer panic(commandError.Error())
		err = gjp.NewJobError(commandError)
		return
	}
	fmt.Printf("The date is %s\n", out)
	return
}

func main() {
	jobPool := gjp.New(2)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Starting jobs again")
		jobPool.Start()
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		jobPool.ShutdownWorkPool()
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		list := jobPool.ListWaitingJobs()
		fmt.Fprintf(w, "%v", list)
	})

	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {

		newJob := &MyJob{
			commandName: "echo",
		}
		_, jobId := jobPool.QueueJob(newJob.commandName, newJob, 0)
		fmt.Fprintf(w, "%v", jobId)
	})

	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) {
		newJob2 := &MyJob{
			commandName: "wrong",
		}
		_, jobId := jobPool.QueueJob(newJob2.commandName, newJob2, 1)
		fmt.Fprintf(w, "%v", jobId)
	})

	http.HandleFunc("/3", func(w http.ResponseWriter, r *http.Request) {

		newJob3 := &MyJob{
			commandName: "test",
		}
		_, jobId := jobPool.QueueJob(newJob3.commandName, newJob3, 1)
		fmt.Fprintf(w, "%v", jobId)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q["id"] != nil {
			searchParam := q["id"][0]
			j, err := jobPool.GetJobFromJobId(searchParam)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				jobson := j.GetJobInfos()
				fmt.Fprintf(w, "%v", string(jobson))
			}
		} else {
			fmt.Fprintf(w, "No search request")
		}
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
