package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jcuga/golongpoll"
)

var savWatcher *fsnotify.Watcher
var debounceTimer *time.Timer
var longPoll *golongpoll.LongpollManager

const debounceInterval = 300 * time.Millisecond

func main() {
	// Set up file watcher and go func to handle events
	var err error
	savWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer savWatcher.Close()
	debounceTimer = time.NewTimer(debounceInterval)
	go watchSavs()

	// TODO: UI for adding watches and/or registry query for SAV locations
	err = savWatcher.Add("F:\\SteamLibrary\\steamapps\\common\\Sid Meier's Civilization III Complete\\Conquests\\Saves")
	if err != nil {
		log.Fatal(err)
	}
	err = savWatcher.Add("F:\\SteamLibrary\\steamapps\\common\\Sid Meier's Civilization III Complete\\Conquests\\Saves\\Auto")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize long poll manager
	longPoll, err = golongpoll.StartLongpoll(golongpoll.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer longPoll.Shutdown()

	// temp loading a hard-coded SAV on startup
	f("F:\\SteamLibrary\\steamapps\\common\\Sid Meier's Civilization III Complete\\Conquests\\Saves\\WAR-Russia-Galley-Mao of the Chinese, 300 AD.SAV")

	// Api server
	server()
}
