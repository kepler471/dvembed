package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

type convertedMedia struct {
	rawMedia
	Ext string
}



// format handles file conversion and formatting with ffmpeg.
// Discord allows a maximum 8MB file upload, so this is only
// required if the size limit is exceeded.
func format(f rawMedia) (*convertedMedia, error) {
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
	var cm convertedMedia
	cm.Path = path.Join(cmd.Dir, f.Id+convertedExt)
	cm.rawMedia.FileInfo, err = os.Stat(cm.Path)
	cm.
	if err != nil {
		log.Println("Error finding formatted file")
	}
	// TODO want to return output error from ffmpeg
	return f, nil
}
