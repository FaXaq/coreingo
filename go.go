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
		commandName string
	}
)

func (myjob *MyJob) ExecuteJob() (err *gjp.JobError) {
	out, error := exec.Command("echo","toto").Output()
	if error != nil {
		defer panic("plz send haelp")
		err = &gjp.JobError{
			error.Error(),
			"nooooooooo",
		}
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

	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {

		newJob := &MyJob{
			commandName: "YES",
		}
		jobPool.QueueJob(newJob.commandName, newJob, 0)
	})

	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) {
		newJob2 := &MyJob{
			commandName: "THAT",
		}
		jobPool.QueueJob(newJob2.commandName, newJob2, 1)
	})

	http.HandleFunc("/3", func(w http.ResponseWriter, r *http.Request) {

		newJob3 := &MyJob{
			commandName: "ROCKS",
		}
		jobPool.QueueJob(newJob3.commandName, newJob3, 1)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
