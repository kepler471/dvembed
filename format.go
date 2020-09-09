package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

type convertedMedia struct {
	os.FileInfo
	Id   string
	Path string
}

// format handles file conversion and formatting with ffmpeg.
// Discord allows a maximum 8MB file upload, so this is only
// required if the size limit is exceeded.
func format(f media) (*convertedMedia, error) {
	cmd := exec.Command(
		"cmd",
		"-i "+f.Name(),
		" -c:v libvpx ",
		"-crf 0 ",
		"-b:v 1M ",
		"-c:a libvorbis ",
		f.Id+convertedExt,
	)
	cmd.Dir = path.Join(dir, f.Id)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed: %v\n", cmd.Args)
	}
	cm := convertedMedia{
		FileInfo: nil,
		Id:       f.Id,
		Path:     "",
	}
	cm.FileInfo, err = os.Stat(cm.Path)
	if err != nil {
		log.Print("Error finding formatted file")
	}
	// TODO want to return output error from ffmpeg
	return &cm, nil
}
