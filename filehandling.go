package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func GetInfosOnMediaFile(fromFile string, infos []string) (info string) {
	var mu sync.Mutex
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd, args := CreateGetInfoCommand(fromFile, infos)

	mu.Lock()
	cmdResult := exec.Command(cmd, args...)
	mu.Unlock()

	cmdResult.Stdout = &out
	cmdResult.Stderr = &stderr

	cmderr := cmdResult.Run()

	if cmderr != nil {
		fmt.Println(cmderr)
		return
	}

	info = strings.TrimRight(out.String(), "\n") //remove the "\n" after the command return
	return
}

func GetFileDuration(file string) (duration int64, err error) {
	lengthString := GetInfosOnMediaFile(file, []string{"duration"})
	duration = 0.0

	if len(lengthString) > 0 {
		durationFloat, errconv := strconv.ParseFloat(lengthString, 64)
		err = errconv
		duration = int64(durationFloat * 1000000) //convert seconds to microseconds
	} else {
		duration = 0
	}

	return
}

//split media files into 5 distinct files
func SplitMediaFile(jobId, file string, duration int64) (files []string, err error) {
	fmt.Println("SplitMediaFile")
	var (
		mu              sync.Mutex
		out             bytes.Buffer
		stderr          bytes.Buffer
		segmentLength   float64
		segmentDuration int
		path            string
		fileExt         string
		fileName        string
		logFileName     string
	)

	//retrieve file informations
	path = GetFileDirectory(file) + "/"
	fileName = GetFileName(file)
	fileExt = GetFileExt(file)
	logFileName = jobId + ".ffconcat"

	if duration%5 == 0 {
		segmentLength = (float64(duration) / 5.0)
	} else {
		segmentLength = (float64(duration) / 5.0) + 1.0 //add the rest of the division to segment
		//to avoid cutting the video
	}

	segmentLength /= 1000000.0
	segmentDuration = RoundUp(segmentLength)

	cmd, args := CreateFileSplitCommand(fileName, path, fileExt, logFileName, segmentDuration)

	mu.Lock()
	cmdResult := exec.Command(cmd, args...)
	mu.Unlock()

	cmdResult.Stdout = &out
	cmdResult.Stderr = &stderr

	cmderr := cmdResult.Run()

	if cmderr != nil {
		fmt.Println(stderr)
		return
	}

	files, err = GetInfosFromFile(
		WorkPath+"/"+logFileName,
		"file", " ", "\n")

	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

func ReplaceFromExtByToExt(id, WorkPath, fromExt, toExt string) (err error) {

	input, err := ioutil.ReadFile(WorkPath + "/" + id + ".ffconcat") //open list

	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, fromExt) {
			lines[i] = strings.Replace(lines[i], fromExt, toExt, -1) //replace all occurences of old ext by new ones
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(WorkPath+"/"+id+".ffconcat", []byte(output), 0644)
	if err != nil {
		return err
	}

	return
}

func GetFileName(file string) (name string) {
	directory := filepath.Dir(file)
	ext := filepath.Ext(file)
	name = file[len(directory)+1 : len(file)-len(ext)] //remove directory and ext from filePath

	return
}

func GetFileExt(file string) (ext string) {
	ext = filepath.Ext(file)
	return
}

func GetFileDirectory(filePath string) (directory string) {
	directory = filepath.Dir(filePath)
	return
}
