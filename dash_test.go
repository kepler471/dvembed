package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_decode(t *testing.T) {
	dashDir := "./DashFiles/"
	fs, _ := ioutil.ReadDir(dashDir)
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".mpd") {
			mpd, _ := os.Open(dashDir + f.Name())
			list := decode(mpd)
			fmt.Println(list)
		}
	}
}

func Test_scrapeReddit(t *testing.T) {
	Url, _ := url.Parse("https://www.reddit.com/domain/v.redd.it/")
	j, _ := fetchJson(Url)
	var result interface{}
	json.Unmarshal(j, &result)
	m := result.(map[string]interface{})

	var stack []string
	var walk func(value reflect.Value)
	walk = func(v reflect.Value) {
		if v.String() == "dash_url" {
			fmt.Printf("Visiting %v\n", v)
		}
		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i))
			}
		case reflect.Map:
			for _, k := range v.MapKeys() {
				//if k.String() == "dash_url" {
				if k.String() == "url_overridden_by_dest" {
					stack = append(stack, fmt.Sprintf("%v", v.MapIndex(k)))
					return
				}
				if k.String() != "secure_media" {
					walk(v.MapIndex(k))
				}
			}
		}
	}
	walk(reflect.ValueOf(m))
	ch := make(chan struct{})
	for _, u := range stack {
		go func(mpd string) {
			resp, _ := fetch(mpd)
			fmt.Println(mpd, decode(resp.Body))
			ch <- struct{}{}
		}(u + "/DASHPlaylist.mpd")
	}
	for range stack {
		<-ch
	}
}
