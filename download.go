package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

type rawMedia struct {
	os.FileInfo
	Id         string
	Url        string
	Path       string
	Downloaded bool
}

// download handles retrieving the v.redd.it media. Initialises
// a rawMedia instance for each URL.
func download(URL string) (*rawMedia, error) {
	f := rawMedia{
		Url: URL,
		Id:  path.Base(URL),
	}
	err := os.MkdirAll(path.Join(dir, f.Id), 0755)
	if err != nil {
		log.Printf("Error creating sub-directory: %v\n", err)
	}
	cmd := exec.Command(
		"youtube-dl",
		"-v",
		"--id",              // use id as name
		"--write-info-json", // save file information
		"--restrict-filenames",
		"--merge-output-format", // Downloading the best available audio and video
		outputFormat,
		URL,
	)
	cmd.Dir = path.Join(dir, f.Id)

	xerr := cmd.Run()
	if xerr != nil {
		log.Fatalf("Failed: %v\n", cmd.Args)
	}

	f.Path = path.Join(cmd.Dir, f.Id+originalExt)
	f.FileInfo, err = os.Stat(f.Path)
	if err != nil {
		log.Println("Error finding downloaded file")
	}
	// TODO want to return output error from youtube-dl
	return &f, xerr
}
