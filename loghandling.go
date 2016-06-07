package main

import (
	"os"
	"fmt"
	"strings"
	"errors"
)

func GetInfosFromFile (filename, infoName, delimiter, dataDelimiter string) (infos []string, err error) {
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

	for i := len(logsArray)-1; i >= 0; i-- {
		if strings.Contains(logsArray[i], infoName + delimiter) {
			infos = append(infos, strings.Trim(logsArray[i],
				infoName + delimiter))
		}
	}

	return
}
