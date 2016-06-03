package main

import (
	"fmt" //format for console logs
	"github.com/FaXaq/gjp" //custom library for jobs
	"os/exec" //command execs
	"net/http" //web services
	"sync" //command sync
	"bytes" //bytes manipulation
	_ "net/http/pprof" //get http logs
	"os" //to create directory
	"strconv" //str conversion
	"errors"
)

type (
	MyJob struct {
		Id string
		Name string
		Command string
		Args []string
		Path string
		MediaLength int64
		MediaFiles []string
		FromFile string
		ToFile string
	}
)

func NewJob (id, command, fromFile, toFile string) (j *MyJob, err error) {
	var (
		cmd string
		args []string
	)
	path := GetFileDirectory(fromFile)
	mediaLength, err := GetFileDuration(fromFile)

	if err != nil {
		fmt.Println("Couldn't get the file duration, set it to 0.0. Error : ",
			err.Error())
	} else {
		fmt.Println("File duration is : ", mediaLength, "µs")
	}

	//create log, tmp, and output directory
	os.Mkdir(path + "/out", 0777)
	os.Mkdir(path + "/log", 0777)
	os.Mkdir(path + "/tmp", 0777)

	mediaFiles, splitErr := SplitMediaFile(fromFile, mediaLength)

	if (splitErr != nil) {
		fmt.Println("Error during split : ", splitErr)
	}

	/* Ends here */

	if command == "convert" {
		cmd, args = CreateConvertCommand(id, path, fromFile, toFile)
	} else if command == "extract-audio" {
		cmd, args = CreateExtractAudioCommand(id, path, fromFile, toFile)
	} else {
		err = errors.New("couldn't create the job, couldn't understand the query")
		return
	}

	j = &MyJob{
		id,
		fromFile + " to " + toFile,
		cmd,
		args,
		path,
		mediaLength,
		mediaFiles,
		fromFile,
		toFile,
	}

	return
}

func (myjob *MyJob) GetProgress(id string) (percentage float64, err error) {
	timings, err := GetInfosFromFile(
		myjob.Path + "/log/" + id + "-logs.gjp",
		"out_time_ms", "=", "\n")

	timingString := timings[len(timings) - 1]

	if len(timingString) > 0 {
		percentage, err = strconv.ParseFloat(timingString, 64) //get the ms timing
		percentage /= float64(myjob.MediaLength) //divide by media length
	} else {
		err = errors.New("No log file or capability to find progress")
	}

	return
}

func (myjob *MyJob) NotifyEnd(id string) {
	http.Get("http://127.0.0.1:8124/")
}

func (myjob *MyJob) NotifyStart(id string) {
	http.Get("http://127.0.0.1:8124/")
}

func (myjob *MyJob) ExecuteJob(id string) (err *gjp.JobError) {
	var mu sync.Mutex
	var out bytes.Buffer
	var stderr bytes.Buffer

	fmt.Println("command : ",myjob.Command, myjob.Args, "on file", myjob.FromFile,"of length", myjob.MediaLength)

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

func ExecutePartialJob (command, file string) (executed bool, err error) {
	return
}
