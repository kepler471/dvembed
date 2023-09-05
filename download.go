package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
//
//	and move "fallback_url" somewhere if its still needed
func fetchJson(Url *url.URL) (*http.Response, error) {
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

	return resp, nil
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
