package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceInterval = 300 * time.Millisecond

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
		}
	}
}

func ExampleNewWatcher() {
	var fn string
	watcher, err := fsnotify.NewWatcher()
	timer := time.NewTimer(debounceInterval)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fn = event.Name
				if event.Op&fsnotify.Write == fsnotify.Write {
					timer.Reset(debounceInterval)
				}
			case <-timer.C:
				// This will get called once debounceInterval after program start, and I'm going to live with that
				f(fn)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("F:\\SteamLibrary\\steamapps\\common\\Sid Meier's Civilization III Complete\\Conquests\\Saves")
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add("F:\\SteamLibrary\\steamapps\\common\\Sid Meier's Civilization III Complete\\Conquests\\Saves\\Auto")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
