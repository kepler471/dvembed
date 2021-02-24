package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_decode(t *testing.T) {
	dashDir := "./DashFiles/"
	fs, _ := ioutil.ReadDir(dashDir)
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".mpd") {
			list := decode(dashDir + f.Name())
			fmt.Println(list)
		}
	}
}
