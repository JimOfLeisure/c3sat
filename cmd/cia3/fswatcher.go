package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func f(s string) {
	fmt.Println(time.Now().String() + " " + s + " modified")
}

func debounce(interval time.Duration, input chan string, f func(arg string)) {
	var (
		filename string
	)
	for {
		select {
		case filename = <-input:
			// fmt.Println("received a send on a spammy channel - might be doing a costly operation if not for debounce")
			fmt.Println("doot")
			// do nothing
		case <-time.After(interval):
			f(filename)
		}
	}
}

func ExampleNewWatcher() {
	watcher, err := fsnotify.NewWatcher()
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
				fn := event.Name
				fi, err := os.Stat(fn)
				if err != nil {
					// log.Fatal(err)
					fmt.Println("stat error - " + fn)
					// return
				} else {
					// fmt.Println(strings.ToLower(fn[len(fn)-4:]))
					if event.Op&fsnotify.Write == fsnotify.Write && strings.ToLower(fn[len(fn)-4:]) == ".sav" && fi.Mode().IsRegular() {
						// log.Println("modified file:", event.Name)
						f(event.Name)
					}
				}
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
