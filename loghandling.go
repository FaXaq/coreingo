package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetInfosFromFile(filename, infoName, delimiter, dataDelimiter string) (infos []string, err error) {
	var start int64

	file, err := os.Open(filename) // For read access
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 500)
	stat, err := os.Stat(filename)

	start = stat.Size() - 500
	if stat.Size() > 500 {
		_, err = file.ReadAt(buf, start)
	} else if stat.Size() > 0 {
		_, err = file.Read(buf)
	} else {
		err = errors.New("Nothing in logfile")
	}

	if err != nil {
		return
	}

	lastInput := string(buf)

	file.Close()

	logsArray := strings.Split(lastInput, dataDelimiter)

	for i := len(logsArray) - 1; i >= 0; i-- {
		if strings.Contains(logsArray[i], infoName+delimiter) {
			infos = append(infos, logsArray[i][len(infoName)+len(delimiter):])
		}
	}

	return
}

func GenerateLogFile(logFileName, fileName, fileExt string) {
	file, err := os.Create(logFileName)

	if err != nil {
		fmt.Println("Error while create logfile")
	}

	file.Write([]byte("ffconcat version 1.0\n"))
	for i := 0; i < 5; i++ {
		file.Write([]byte("file " + fileName + "-" + strconv.Itoa(i) + fileExt + "\n"))
	}
}
