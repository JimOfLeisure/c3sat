package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/myjimnelson/c3sat/civ3satgql"
)

func f(s string) {
	// fmt.Println(s + " <==")
	if len(s) > 4 && strings.ToLower(s[len(s)-4:]) == ".sav" {
		fi, err := os.Stat(s)
		if err != nil {
			// log.Fatal(err)
			fmt.Println("stat error - " + s)
			return
		}
		if fi.Mode().IsRegular() {
			fmt.Println(time.Now().String() + " " + s + " modified")
			err := civ3satgql.ChangeSavePath(s)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func watchSavs() {
	var fn string
	for {
		select {
		case event, ok := <-savWatcher.Events:
			if !ok {
				return
			}
			fn = event.Name
			if event.Op&fsnotify.Write == fsnotify.Write {
				debounceTimer.Reset(debounceInterval)
			}
		case <-debounceTimer.C:
			// This will get called once debounceInterval after program start, and I'm going to live with that
			f(fn)
		case err, ok := <-savWatcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
