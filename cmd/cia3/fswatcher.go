package main

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/myjimnelson/c3sat/civ3satgql"
)

type watchListType struct {
	watches []string
}

// Does not check to see if already added
func (w *watchListType) addWatch(path string) error {
	err := savWatcher.Add(path)
	if err != nil {
		return err
	}
	w.watches = append(w.watches, path)
	return nil
}

// Only deletes one
func (w *watchListType) removeWatch(path string) error {
	err := savWatcher.Remove(path)
	if err != nil {
		return err
	}
	for i := 0; i < len(w.watches); i++ {
		if w.watches[i] == path {
			// remove element from array by swapping last element and replacing with one-shorter array
			w.watches[i] = w.watches[len(w.watches)-1]
			w.watches[len(w.watches)-1] = ""
			w.watches = w.watches[:len(w.watches)-1]
			break
		}
	}
	return nil
}

func loadDefaultBiq(s string) error {
	fi, err := os.Stat(s)
	if err != nil {
		return err
	}
	if fi.Mode().IsRegular() {
		err := civ3satgql.ChangeDefaultBicPath(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadNewSav(s string) error {
	if len(s) > 4 && strings.ToLower(s[len(s)-4:]) == ".sav" {
		fi, err := os.Stat(s)
		if err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			err := civ3satgql.ChangeSavePath(s)
			if err != nil {
				return err
			}
			longPoll.Publish("refresh", s)
		}
	}
	return nil
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
