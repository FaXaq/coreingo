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
	}

	if GetFileExt(toFile) == ".flv" || GetFileExt(fromFile) == ".flv" {
		args = append(args, "-c:v")
		args = append(args, "libx264")
	}

	args = append(args,
		WorkPath + "/" + toFile)

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

func CreateFileSplitCommand(fileName, fromExt, path, toFileName, toExt, logFileName string, duration int) (cmd string, args []string) {
	cmd = "ffmpeg"
	args = []string{
		"-y",
		"-i",
		path + fileName + fromExt,
	}

	for i := 0; i < 5; i++ {
		args = append(args,
			[]string{
				"-codec",
				"copy",
				"-ss",
				strconv.Itoa(duration * i),
				"-t",
				strconv.Itoa(duration),
				WorkPath + "/" + toFileName + "-" + strconv.Itoa(i) + toExt,
			}...)
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
