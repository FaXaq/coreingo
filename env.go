package main

import (
	"errors"
	"os"
	"path/filepath"
	"fmt"
)

func GetEnvConfig() (err error) {
	fmt.Pritln("=== Config ===")

	CallbackStart = os.Getenv("GJP_STARTJOB_CALLBACK")
	if CallbackStart == "" {
		err = errors.New("Couldn't get GJP_STARTJOB_CALLBACK env variable")
		return
	} else {
		fmt.Println("Callback start :", CallbackStart)
	}

	CallbackEnd = os.Getenv("GJP_ENDJOB_CALLBACK")
	if CallbackEnd == "" {
		err = errors.New("Couldn't get GJP_ENDJOB_CALLBACK env variable")
		return
	} else {
		fmt.Println("Callback end :", CallbackEnd)
	}

	LogPath = filepath.Dir(os.Getenv("GJP_LOG_PATH"))
	if LogPath == "" {
		err = errors.New("Couldn't get GJP_LOG_PATH env variable")
		return
	} else {
		fmt.Println("Log path :", LogPath)
	}


	WorkPath = filepath.Dir(os.Getenv("GJP_WORK_PATH"))
	if WorkPath == "" {
		err = errors.New("Couldn't get GJP_WORK_PATH env variable")
		return
	} else {
		fmt.Println("WorkPath :", WorkPath)
	}

	Port = os.Getenv("GJP_PORT")
	if Port == "" {
		err = errors.New("Couldn't get GJP_PORT env variable")
		return
	} else {
		fmt.Println("Port :", Port)
	}

	fmt.Println("=== End of Config ===")

	return
}
