package main

import (
	"fmt" //format for console logs
	"github.com/FaXaq/gjp" //custom library for jobs
	"os/exec" //command execs
	"net/http" //web services
	"html" //html answer
	"log" //log
	"sync" //command sync
	"bytes" //bytes manipulation
	_ "net/http/pprof" //get http logs
	"strconv" //help conversion
	"strings" //manipulate strings
	"path/filepath" //to get filepath from string
	"os" //to create directory
)

type (
	MyJob struct {
		Name string
		Command string
		Args []string
		MediaLength float64
		FromFile string
		ToFile string
	}
)

func NewJob (command, fromFile, toFile string) (j *MyJob) {
	/* To get the initial media length for percentage purposes */

	var mu sync.Mutex
	var out bytes.Buffer
	var stderr bytes.Buffer

	currentPath := filepath.Dir(fromFile)

	//create log & output directory
	os.Mkdir(currentPath + "/out", 0777)
	os.Mkdir(currentPath + "/log", 0777)

	cmd := "ffprobe"
	args := []string{"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=duration",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		fromFile}

	mu.Lock()
	cmdResult := exec.Command(cmd,args...)
	mu.Unlock()

	cmdResult.Stdout = &out //for debug purposes
	cmdResult.Stderr = &stderr //for debug purposes

	cmderr := cmdResult.Run()

	if cmderr != nil {
		fmt.Println(cmderr)
		return
	}

	mediaLength := 0.0

	if out.Len() > 0 {
		mediaLength, _ = strconv.ParseFloat(strings.TrimRight(out.String(), "\n"),64)
	}

	fmt.Println(mediaLength)

	/* Ends here */

	if command == "convert" {
		/* build the command arguments */
		args := []string{
			"-y",
			"-i",
			fromFile,
			"-progress",
			currentPath + "/log/logs.gjp",
			currentPath + "/out" + toFile,
		}
		/* ends here */

		j = &MyJob{
			fromFile + " to " + toFile,
			"ffmpeg",
			args,
			mediaLength,
			fromFile,
			toFile,
		}
	}

	return
}

func GetTimingFromLogFile(filename string) (currentTiming string) {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 200)
	stat, err := os.Stat(filename)
	start := stat.Size() - 200
	_, err = file.ReadAt(buf, start)
	if err == nil {
		fmt.Printf("%s\n", buf)
	}
	file.Close()
	return
}

func (myjob *MyJob) GetProgress(id string) {
	GetTimingFromLogFile("logs.gjp")
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

	//myjob.Args = append(myjob.Args, myjob.FileName) //No need for now

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

	// http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
	// 	list := jobPool.ListWaitingJobs()
	// 	fmt.Fprintf(w, "%v", list)
	// })

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q["command"] != nil && q["fromFile"] != nil && q["toFile"] != nil {
			newJob := NewJob(
				q["command"][0],
				q["fromFile"][0],
				q["toFile"][0])
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

	http.HandleFund("/getProgress", func(w http.ResponseWriter, r *http.Request) {
		_ := r.URL.Query()
	})

	log.Println(http.ListenAndServe("localhost:6060", nil))
}
