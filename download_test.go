package main

import (
	"net/url"
	"testing"
)

func TestDownloadVRedditLink(t *testing.T) {
	URLs := []string{
		"https://v.redd.it/duir5tuwswl51",
		"https://v.redd.it/5ltubsoyawl51",
		"https://v.redd.it/dttgnvp69wl51",
		"https://v.redd.it/e497qwjsh1m51",
		"https://v.redd.it/b8b0ha50l5g41",
		"https://v.redd.it/zv89llsvexdz",
	}
	for _, URL := range URLs {
		Url, _ := url.Parse(URL)
		_, err := download(Url)
		if err != nil {
			t.Errorf(`youtube-dl %q failed`, URL)
		}
	}
}

func TestDownloadRedditLink(t *testing.T) {
	URLs := []string{
		"https://www.reddit.com/r/IdiotsInCars/comments/ioqqbf/i_know_ill_cut_in_front_of_this_semi/",
		"https://www.reddit.com/r/AnimalsBeingBros/comments/ip89wl/possibly_the_most_patient_kitty_in_the_world_with/",
		"https://old.reddit.com/r/StarWars/comments/l1l7f6/finished_this_last_night_took_me_30_hours_to/",
		"https://www.reddit.com/r/videos/comments/6rrwyj/that_small_heart_attack/",
		"https://www.reddit.com/r/videos/comments/6rrwyj",
		"https://www.reddit.com/r/MadeMeSmile/comments/6t7wi5/wait_for_it/",
		"https://old.reddit.com/r/MadeMeSmile/comments/6t7wi5/wait_for_it/",
		"https://www.reddit.com/r/videos/comments/6t7sg9/comedians_hilarious_joke_about_the_guam_flag/",
		"https://www.reddit.com/r/videos/comments/6t75wq/southern_man_tries_to_speak_without_an_accent/",
		"https://nm.reddit.com/r/Cricket/comments/8idvby/lousy_cameraman_finds_himself_in_cairns_line_of/",
	}
	for _, URL := range URLs {
		Url, _ := url.Parse(URL)
		_, err := download(Url)
		if err != nil {
			t.Errorf(`youtube-dl %q failed`, URL)
			// The following logs were producing some errors, have been removed
			//t.Log(err)
			//t.Logf("Id:\t%v,\nPath:\t%v,\nName:\t%v,\nSize:\t%v", f.Id, f.Path, f.Name(), f.Size())
			// Output:
			/*
				=== RUN   TestDownloadVRedditLink
				2020/09/09 19:18:57 ...run youtube-dl...
				2020/09/09 19:18:58 ...run youtube-dl...
				2020/09/09 19:18:59 ...run youtube-dl...
				2020/09/09 19:19:00 ...run youtube-dl...
				2020/09/09 19:19:00 Error finding downloaded file: stat downloads/e497qwjsh1m51/e497qwjsh1m51.mp4: no such file or directory
				    download_test.go:15: youtube-dl "https://v.redd.it/e497qwjsh1m51" failed
				--- FAIL: TestDownloadVRedditLink (3.25s)
				=== RUN   TestDownloadRedditLink
				2020/09/09 19:19:00 ...run youtube-dl...
				2020/09/09 19:19:04 ...run youtube-dl...
				2020/09/09 19:19:06 Error finding downloaded file: stat downloads/possibly_the_most_patient_kitty_in_the_world_with/possibly_the_most_patient_kitty_in_the_world_with.mp4: no such file or directory
				    download_test.go:27: youtube-dl "https://www.reddit.com/r/AnimalsBeingBros/comments/ip89wl/possibly_the_most_patient_kitty_in_the_world_with/" failed
				    download_test.go:28: stat downloads/possibly_the_most_patient_kitty_in_the_world_with/possibly_the_most_patient_kitty_in_the_world_with.mp4: no such file or directory
				--- FAIL: TestDownloadRedditLink (5.44s)
				panic: runtime error: invalid memory address or nil pointer dereference [recovered]
				        panic: runtime error: invalid memory address or nil pointer dereference
				[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x6ed596]

				goroutine 7 [running]:
				testing.tRunner.func1.1(0x73b8a0, 0xa459e0)
				        /usr/lib/golang/src/testing/testing.go:988 +0x30d
				testing.tRunner.func1(0xc00011c900)
				        /usr/lib/golang/src/testing/testing.go:991 +0x3f9
				panic(0x73b8a0, 0xa459e0)
				        /usr/lib/golang/src/runtime/panic.go:969 +0x166
				dvembed.TestDownloadRedditLink(0xc00011c900)
				        /home/steli/source/dvembed/download_test.go:29 +0x1f6
				testing.tRunner(0xc00011c900, 0x7be460)
				        /usr/lib/golang/src/testing/testing.go:1039 +0xdc
				created by testing.(*T).Run
				        /usr/lib/golang/src/testing/testing.go:1090 +0x372
				exit status 2
				FAIL    dvembed 8.691s

			*/
		}
	}
}
