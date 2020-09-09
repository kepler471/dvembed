package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

// Format handles file conversion and formatting with ffmpeg.
// Discord allows a maximum 8MB file upload, so this is only
// required if the size limit is exceeded.
func Format(m Media) (Media, error) {
	cmd := exec.Command(
		"cmd",
		"-i "+m.Info.Name(),
		" -c:v libvpx ",
		"-crf 0 ",
		"-b:v 1M ",
		"-c:a libvorbis ",
		m.Id+ConvertedExt,
	)
	cmd.Dir = path.Join(Dir, m.Id)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed: %v\n", cmd.Args)
	}
	m.Info, err = os.Stat(path.Join(cmd.Dir, m.Id+ConvertedExt))
	return m, err
}
