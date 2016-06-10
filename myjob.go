package main

import (
	"bytes" //bytes manipulation
	"errors"
	"fmt"                  //format for console logs
	"github.com/FaXaq/gjp" //custom library for jobs
	"io/ioutil"
	"net/http"
	"os"      //to create directory
	"os/exec" //command execs
	"strconv" //str conversion
	"sync"    //command sync
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
		Splitted      bool
	}
)

func NewJob(id, command, fromFile, toFile string) (j *MyJob, err error) {
	var (
		cmd        string     //tmp command
		cmds       []string   //commands array
		args       []string   //args array
		cmdsArgs   [][]string //commands arguments matching commands array
		toExt      string
		path       string //working directory
		splitted   bool
		mediaFiles []string
	)

	path = GetFileDirectory(fromFile)
	mediaLength, err := GetFileDuration(fromFile)
	toExt = GetFileExt(toFile)

	if err != nil {
		fmt.Println("Couldn't get the file duration, set it to 0.0. Error : ",
			err.Error())
		splitted = false
	} else {
		fmt.Println("File duration is : ", mediaLength, "µs")
		splitted = true
	}

	mediaFiles, splitErr := SplitMediaFile(id, fromFile, mediaLength)

	if splitErr != nil {
		fmt.Println("Error during split : ", splitErr)
	}

	if len(mediaFiles) > 1 {
		ReplaceFromExtByToExt(id, toExt, mediaFiles)
	}

	fmt.Println(mediaFiles)

	//create log, tmp, and output directory
	os.Mkdir(path+"/"+"out", 0777) //filepath.Div
	os.Mkdir(LogPath+"/"+id, 0777)
	os.Mkdir(WorkPath, 0777)

	if command == "convert" {
		for i := 0; i < len(mediaFiles); i++ {
			toPartFile := GetFileName(mediaFiles[i]) + toExt
			cmd, args = CreateConvertCommand(GetFileName(mediaFiles[i]),
				path,           //fromFile
				LogPath+"/"+id, //jobLogPath
				mediaFiles[i],
				toPartFile)
			cmds = append(cmds, cmd)
			cmdsArgs = append(cmdsArgs, args)
		}
	} else if command == "extract-audio" {
		for i := 0; i < len(mediaFiles); i++ {
			toPartFile := GetFileName(mediaFiles[i]) + toExt
			cmd, args = CreateExtractAudioCommand(GetFileName(mediaFiles[i]),
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
		splitted,
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
	fmt.Println("Notifying start of job ", j.GetJobId(), ":>", CallbackStart)

	jobson, _ := j.GetJobInfos()

	fmt.Println(string(jobson))

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
		fmt.Println("Couldn't notify start job")
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

	if myjob.Splitted == false {
		return
	}

	//concat file when finished
	command, args := CreateConcatCommand(WorkPath+"/"+j.GetJobId()+".ffconcat", myjob.Path, myjob.ToFile)

	fmt.Println(command, args)

	mu.Lock()
	out, cmderr := exec.Command(command, args...).Output()
	mu.Unlock()

	if cmderr != nil {
		err = gjp.NewJobError(cmderr, string(out))
	}

	//retrieve new file names and add them to media files of the job
	newMediaFiles, newMediaFilesErr :=
		GetInfosFromFile(WorkPath+"/"+j.GetJobId()+".ffconcat", "file", " ", "\n")

	if newMediaFilesErr != nil {
		err = gjp.NewJobError(newMediaFilesErr, newMediaFilesErr.Error())
	}

	for _, newMediaFile := range newMediaFiles {
		myjob.MediaFiles = append(myjob.MediaFiles, newMediaFile)
	}

	//remove all tmp files
	removeErr := RemoveFilesFromWorkDir(append(myjob.MediaFiles, j.GetJobId()+".ffconcat"))

	if removeErr != nil {
		err = gjp.NewJobError(removeErr, removeErr.Error())
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
