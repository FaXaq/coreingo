package main

import (
	"bytes" //bytes manipulation
	"errors"
	"fmt"                  //format for console logs
	"github.com/FaXaq/gjp" //custom library for jobs
	"io/ioutil"
	"net/http"
	"os"               //to create directory
	"os/exec"          //command execs
	"strconv"          //str conversion
	"sync"             //command sync
)

type (
	MyJob struct {
		Id            string
		Name          string
		Commands      []string
		Args          [][]string
		Path          string
		MediaLength   int64
		MediaFiles    []string
		FromFile      string
		ToFile        string
		ReportChannel chan *gjp.JobError
	}
)

func NewJob(id, command, fromFile, toFile string) (j *MyJob, err error) {
	var (
		cmd          string     //tmp command
		cmds         []string   //commands array
		args         []string   //args array
		cmdsArgs     [][]string //commands arguments matching commands array
		fromFileName string
		toExt        string
		fromExt      string
		path         string //working directory
	)
	path = GetFileDirectory(fromFile)
	mediaLength, err := GetFileDuration(fromFile)
	fromFileName = GetFileName(fromFile)
	fromExt = GetFileExt(fromFile)
	toExt = GetFileExt(toFile)

	if err != nil {
		fmt.Println("Couldn't get the file duration, set it to 0.0. Error : ",
			err.Error())
	} else {
		fmt.Println("File duration is : ", mediaLength, "µs")
	}

	//create log, tmp, and output directory
	os.Mkdir(GetFileDirectory(path)+"/"+"out", 0777) //filepath.Div
	os.Mkdir(LogPath+"/"+id, 0777)
	os.Mkdir(WorkPath, 0777)

	mediaFiles, splitErr := SplitMediaFile(id, fromFile, mediaLength)
	ReplaceFromExtByToExt(id, WorkPath, fromExt, toExt)

	if splitErr != nil {
		fmt.Println("Error during split : ", splitErr)
	}

	if command == "convert" {
		for i := 0; i < len(mediaFiles); i++ {
			toPartFile := fromFileName + "-" + strconv.Itoa(i) + toExt
			cmd, args = CreateConvertCommand(id+"-"+strconv.Itoa(i),
				path,           //fromFile
				LogPath+"/"+id, //jobLogPath
				mediaFiles[i],
				toPartFile)
			cmds = append(cmds, cmd)
			cmdsArgs = append(cmdsArgs, args)
		}
	} else if command == "extract-audio" {
		for i := 0; i < len(mediaFiles); i++ {
			toPartFile := fromFileName + "-" + strconv.Itoa(i) + toExt
			cmd, args = CreateExtractAudioCommand(id+"-"+strconv.Itoa(i),
				path,           //fromFile
				LogPath+"/"+id, //jobLogPath
				mediaFiles[i],
				toPartFile)
			cmds = append(cmds, cmd)
			cmdsArgs = append(cmdsArgs, args)
		}
	} else {
		err = errors.New("couldn't create the job, couldn't understand the query")
		return
	}

	j = &MyJob{
		id,
		fromFile + " to " + toFile,
		cmds,
		cmdsArgs,
		path,
		mediaLength,
		mediaFiles,
		fromFile,
		toFile,
		make(chan *gjp.JobError, 2),
	}

	return
}

func (myjob *MyJob) GetProgress(id string) (percentage float64, err error) {
	var (
		timings   []string
		timing    string
		timingSum float64
	)
	jobLogPath := LogPath + "/" + id + "/"
	files, err := ioutil.ReadDir(jobLogPath)

	timingSum = 0.0

	//retrieve all timings from files
	for _, file := range files {
		timingInfos, err := GetInfosFromFile(jobLogPath+file.Name(),
			"out_time_ms", "=", "\n")

		timing = timingInfos[len(timingInfos)-1]

		timings = append(timings, timing)

		if err != nil {
			fmt.Println(err.Error())
			return timingSum, err
		}
	}

	for i := 0; i < len(timings); i++ {
		tmpTiming, err := strconv.ParseFloat(timings[len(timings)-1], 64)
		if err != nil {
			fmt.Println(err.Error())
			return timingSum, err
		}
		timingSum += tmpTiming
	}

	//get the ms timing
	percentage = timingSum / float64(myjob.MediaLength) //divide by media length

	return
}

func (myjob *MyJob) NotifyEnd(j *gjp.Job) {
	fmt.Println("Notifying end of job ", j.GetJobId(), ":>", CallbackEnd)

	jobson, _ := j.GetJobInfos()
	req, err := http.NewRequest("POST", CallbackEnd, bytes.NewBuffer(jobson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error while getting response from NotifyEnd",
			err.Error())
		return
	}
	defer resp.Body.Close()


	if resp.Status == "200" {
		fmt.Println("Notified end job")
	} else {
		fmt.Println("Couldn't notify end job")
	}
}

func (myjob *MyJob) NotifyStart(j *gjp.Job) {
	fmt.Println("Notifying end of job ", j.GetJobId(), ":>", CallbackStart)

	jobson, _ := j.GetJobInfos()
	req, err := http.NewRequest("POST", CallbackStart, bytes.NewBuffer(jobson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error while getting response from NotifyEnd",
			err.Error())
		return
	}
	defer resp.Body.Close()


	if resp.Status == "200" {
		fmt.Println("Notified end job")
	} else {
		fmt.Println("Couldn't notify end job")
	}
}

func (myjob *MyJob) ExecuteJob(j *gjp.Job) (err *gjp.JobError) {
	var (
		mu  sync.Mutex
		out []byte
	)

	for i := 0; i < len(myjob.Commands); i++ {
		go myjob.ExecutePartialJob(myjob.Commands[i],
			myjob.Args[i])
	}

	for i := 0; i < len(myjob.Commands); i++ {
		err = <-myjob.ReportChannel
		if err != nil {
			fmt.Println("\n------\nPart", i, "of Job", j.GetJobId(), "errored\n------\n")
			fmt.Println(err.FmtError())
		}
	}

	fmt.Println("finished job, now concat")

	//concat file when finished

	command, args := CreateConcatCommand(WorkPath + "/" + j.GetJobId() + ".ffconcat", myjob.Path, myjob.ToFile)

	fmt.Println(command, args)

	mu.Lock()
	out, cmderr := exec.Command(command, args...).Output()
	mu.Unlock()

	if cmderr != nil {
		err = gjp.NewJobError(cmderr, string(out))
	}

	return
}

func (myjob *MyJob) ExecutePartialJob(command string, args []string) {
	var (
		mu     sync.Mutex
		out    bytes.Buffer
		stderr bytes.Buffer
		err    *gjp.JobError
	)

	fmt.Println(command, args)

	mu.Lock()
	cmd := exec.Command(command, args...)
	mu.Unlock()

	cmd.Stdout = &out    //for debug purposes
	cmd.Stderr = &stderr //for debug purposes

	cmderr := cmd.Run()

	if cmderr != nil {
		err = gjp.NewJobError(cmderr, stderr.String())
	}

	myjob.ReportChannel <- err

	return
}
