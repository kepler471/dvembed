package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
)

type media struct {
	os.FileInfo
	Id   string
	Path string
}

// download handles retrieving the v.redd.it media using youtube-dl.
// Initialises a media instance for each URL.
func download(Url *url.URL) (*media, error) {
	f := media{
		Id: path.Base(Url.String()),
	}
	cmd := exec.Command(
		"youtube-dl",
		"-v",
		//"--id",              // use id as name
		"--output",
		f.Id+originalExt,
		"--write-info-json", // save file information
		//"--restrict-filenames",
		"--merge-output-format", // Downloading the best available audio and video
		outputFormat,
		Url.String(),
	)
	cmd.Dir = path.Join(dir, f.Id)
	err := os.MkdirAll(cmd.Dir, 0755)
	if err != nil {
		log.Printf("\tError creating sub-directory: %v", err)
		cmd.Dir = dir
	}

	f.Path = path.Join(cmd.Dir, f.Id+originalExt)
	if info, err := os.Stat(f.Path); !os.IsNotExist(err) {
		log.Printf("\tFile at %v exists, will not be downloaded. Err: %v", f.Path, err)
		f.FileInfo = info
		return &f, nil
	} else {
		log.Print("\t...run youtube-dl...")
		err = cmd.Run()
		if err != nil {
			log.Printf("\tFailed youtube-dl process: %v", cmd.Args)
			return &f, err
		}
	}
	f.FileInfo, err = os.Stat(f.Path)
	if err != nil {
		log.Print("\tError finding downloaded file: ", err)
	}
	// TODO want to return output error from youtube-dl
	return &f, err
}
