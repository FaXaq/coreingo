package main

import (
	"errors"
	"os"
	"path/filepath"
)

func GetEnvConfig() (err error) {
	CallbackStart = os.Getenv("GJP_STARTJOB_CALLBACK")
	if CallbackStart == "" {
		err = errors.New("Couldn't get GJP_STARTJOB_CALLBACK env variable")
		return
	}

	CallbackEnd = os.Getenv("GJP_ENDJOB_CALLBACK")
	if CallbackEnd == "" {
		err = errors.New("Couldn't get GJP_ENDJOB_CALLBACK env variable")
		return
	}

	LogPath = filepath.Dir(os.Getenv("GJP_LOG_PATH"))
	if LogPath == "" {
		err = errors.New("Couldn't get GJP_LOG_PATH env variable")
		return
	}

	WorkPath = filepath.Dir(os.Getenv("GJP_WORK_PATH"))
	if WorkPath == "" {
		err = errors.New("Couldn't get GJP_WORK_PATH env variable")
		return
	}

	Port = os.Getenv("GJP_PORT")
	if Port == "" {
		err = errors.New("Couldn't get GJP_PORT env variable")
		return
	}

	return
}
