package main

import (
	"strconv"
	"strings"
)

/*
create commands
*/

func CreateConvertCommand(id, path, logPath, fromFile, toFile string) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-y",
		"-i",
		WorkPath + "/" + fromFile,
		"-strict",
		"-2",
		"-progress",
		logPath + "/" + id + "-logs.gjp",
		WorkPath + "/" + toFile,
	}

	return
}

func CreateExtractAudioCommand(id, path, logPath, fromFile, toFile string) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-y",
		"-i",
		WorkPath + "/" + fromFile,
		"-strict",
		"-2",
		"-progress",
		logPath + "/" + id + "-logs.gjp",
		WorkPath + "/" + toFile,
	}

	return
}

func CreateGetInfoCommand(fromFile string, infos []string) (cmd string, args []string) {
	cmd = "ffprobe"
	args = []string{"-v",
		"error",
		"-show_entries",
		"format=" + strings.Join(infos, ","),
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		fromFile}

	return
}

func CreateFileSplitCommand(fileName, fromExt, path, toExt, logFileName string, duration int) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-i",
		path + fileName + fromExt,
		"-acodec",
		"copy",
		"-f",
		"segment",
		"-segment_time",
		strconv.Itoa(duration),
		"-vcodec",
		"copy",
		"-reset_timestamps",
		"1",
		"-map",
		"0",
		"-segment_list",
		WorkPath + "/" + logFileName,
		"-segment_list_type",
		"ffconcat",
		WorkPath + "/" + fileName + "-%d" + toExt,
	}

	return
}

func CreateConcatCommand(inputFile, path, toFile string) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-y",
		"-f",
		"concat",
		"-i",
		inputFile,
		"-c",
		"copy",
		path + "/out/" + toFile,
	}

	return
}
