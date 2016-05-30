package main

import (
	"fmt"
	"github.com/FaXaq/gjp"
	"os/exec"
	"net/http"
	"html"
	"log"
	"sync"
	"bytes"
	_ "net/http/pprof"
)

type (
	MyJob struct {
		Name string
		Command string
		Args []string
	}
)

func (myjob *MyJob) GetProgress(id string) {
	fmt.Println("100% Maggle")
}

func (myjob *MyJob) NotifyEnd(id string) {
	http.Get("http://127.0.0.1:8124/")
}

func (myjob *MyJob) NotifyStart(id string) {
	http.Get("http://127.0.0.1:8124/")
}

func (myjob *MyJob) ExecuteJob() (err *gjp.JobError) {
	var mu sync.Mutex
	var out bytes.Buffer
	var stderr bytes.Buffer
	fmt.Println("command : ",myjob.Command, myjob.Args)

	mu.Lock()
	cmd := exec.Command(myjob.Command,myjob.Args...)
	mu.Unlock()

	cmd.Stdout = &out //for debug purposes
	cmd.Stderr = &stderr //for debug purposes

	cmderr := cmd.Run()

	if cmderr != nil {
		err = gjp.NewJobError(cmderr, stderr.String())
		return
	}
	fmt.Println(out.String())

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

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q["name"] != nil && q["command"] != nil {
			newJob := &MyJob{
				q["name"][0],
				q["command"][0],
				q["args"],
			}
			_, jobId := jobPool.QueueJob(newJob.Name, newJob, 0)
			fmt.Fprintf(w, "%v", jobId)
		} else {
			fmt.Fprintf(w, "Couldn't create the job")
		}
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

	log.Println(http.ListenAndServe("localhost:6060", nil))
}
