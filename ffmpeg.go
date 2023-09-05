package main

import (
	"fmt"
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

// mux uses ffmpeg to merge the video from one file with the audio from another.
// The command used:
//
//	ffmpeg -y -loglevel repeat+info \
//		-i file:video.mp4 \
//		-i file:audio.mp4 \
//		-c copy -map 0:v:0 -map 1:a:0 file:output.mp4
func mux(video, audio, output, cwd string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-loglevel",
		"repeat+info",
		"-i",
		"file:"+video,
		"-i",
		"file:"+audio,
		"-c",
		"copy",
		"-map",
		"0:v:0",
		"-map",
		"1:a:0",
		"file:"+output,
	)
	cmd.Dir = cwd
	// TODO: log cmd.Stdout to ffmpeg.log file
	cmd.Stderr = os.Stderr
	log.Printf("\t[RUN] %v", cmd.Args)
	log.Printf("%v", cmd.Dir)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed ffmpeg mux: %v", cmd.Stderr)
	}
	return nil
}

// compress handles file conversion and formatting with ffmpeg.
// Discord allows a maximum 8MB file upload, so this is only
// required if the size limit is exceeded.
func compress(f media) (*convertedMedia, error) {
	var crf = 40
	//cmd := exec.Command(
	//	"ffmpeg",
	//	"-i "+f.Name(),
	//	" -c:v libvpx ",
	//	fmt.Sprintf("-crf %v ", crf),
	//	"-b:v 1M ",
	//	"-c:a libvorbis ",
	//	f.Id+convertedExt,
	//)
	cmd := exec.Command(
		"ffmpeg",
		fmt.Sprintf("-i %v -c:v libvpx -crf %v -b:v 1M -c:a libvorbis %v",
			f.Name(), crf, f.Id+convertedExt),
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
