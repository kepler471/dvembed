package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type media struct {
	os.FileInfo
	Id   string
	Path string
}

func fetch(Url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, fmt.Errorf("\t...Failed creating request for %v: %v", Url, err)
	}

	req.Header.Set("User-Agent", "Dvembed/0.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("\t...Failed `GET` request for %v, %v", Url, err)
	}

	return resp, nil
}

// TODO: rename to fetchJson
// 	and move "fallback_url" somewhere if its still needed
func fetchJson(Url *url.URL) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", Url.String()+".json", nil)
	if err != nil {
		return nil, fmt.Errorf("\t...Failed creating request for %v: %v", Url.String()+".json", err)
	}

	req.Header.Set("User-Agent", "Dvembed/0.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("\t...Failed `GET` request for %v, %v", Url.String()+".json", err)
	}

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("\t...Error reading response, %v", err)
	}

	if !bytes.Contains(j, []byte(`"fallback_url"`)) {
		return j, fmt.Errorf("\t...fallback_url not found.")
	}

	return j, nil
}

func download(Url, directory string) (f *os.File, written int64, err error) {
	// see https://gobyexample.com/worker-pools for goroutines and channels
	f, err = os.Create(filepath.Join(directory, path.Base(Url)))
	//defer f.Close()
	if err != nil {
		return nil, 0, fmt.Errorf("\tcould not create file, %v\n", err)
	}
	data, err := fetch(Url)
	if err != nil {
		return nil, 0, fmt.Errorf("\terror fetching file: %v\n", err)
	}
	defer data.Body.Close()
	written, err = io.Copy(f, data.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error writing to %v: %v\n", f.Name(), err)
	}
	return f, written, nil
}

// youtubeDlDownload handles retrieving the v.redd.it media using youtube-dl.
// Initialises a media instance for each URL.
func youtubeDlDownload(Url *url.URL) (*media, error) {
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
