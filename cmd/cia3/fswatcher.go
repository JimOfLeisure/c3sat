package main

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/myjimnelson/c3sat/civ3satgql"
)

func loadNewSav(s string) {
	if len(s) > 4 && strings.ToLower(s[len(s)-4:]) == ".sav" {
		fi, err := os.Stat(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		if fi.Mode().IsRegular() {
			err := civ3satgql.ChangeSavePath(s)
			if err != nil {
				log.Fatal(err)
			}
			longPoll.Publish("refresh", s)
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
			loadNewSav(fn)
		case err, ok := <-savWatcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
