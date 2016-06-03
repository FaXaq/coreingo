package main

import (
	"strings"
	"strconv"
)

/*
create commands
*/

func CreateConvertCommand (id, path, fromFile, toFile string) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-y",
		"-i",
		fromFile,
		"-progress",
		path + "/log/" + id + "-logs.gjp",
		path + "/out/" + toFile,
	}

	return
}

func CreateExtractAudioCommand(id, path, fromFile, toFile string) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-y",
		"-i",
		fromFile,
		"-progress",
		path + "/log/" + id + "-logs.gjp",
		path + "/out/" + toFile,
	}

	return
}

func CreateGetInfoCommand (fromFile string, infos []string) (cmd string, args []string) {
	cmd = "ffprobe"
	args = []string{"-v",
		"error",
		"-show_entries",
		"format=" + strings.Join(infos,","),
		"-of",
		"default=noprint_wrappers=1:nokey=1",
		fromFile}

	return
}

func CreateFileSplitCommand (fileName, path, fileExt, logFileName string, duration int) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-y",
		"-i",
		path + fileName + fileExt,
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
		path + "tmp/" + logFileName,
		"-segment_list_type",
		"ffconcat",
		path + "tmp/" + fileName + "-%d" + fileExt,
	}

	return
}
