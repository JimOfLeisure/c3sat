package main

import (
	// "fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jcuga/golongpoll"
	"github.com/zserge/lorca"
)

var savWatcher *fsnotify.Watcher
var debounceTimer *time.Timer
var longPoll *golongpoll.LongpollManager

const debounceInterval = 300 * time.Millisecond

func main() {
	// fmt.Println("\nCiv Intelligence Agency III alpha 1b\n")
	// fmt.Println("Setting up\n")
	// Set up file watcher
	var err error
	savWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer savWatcher.Close()

	// Set up file event handler
	debounceTimer = time.NewTimer(debounceInterval)
	go watchSavs()

	// Initialize long poll manager
	longPoll, err = golongpoll.StartLongpoll(golongpoll.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer longPoll.Shutdown()

	// Read Win registry for Civ3 Conquests path
	civPath, err := findWinCivInstall()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Detected Civ3 location: " + civPath + "\n")

	lastSav, err := getLastSav(civPath)
	if err != nil {
		// fmt.Println("Failed to discover latest save from conquests.ini. " + err.Error())
	} else {
		// fmt.Println("Opening latest SAV file " + lastSav + "\n")
		loadNewSav(lastSav)
	}

	// fmt.Println(`Adding <civ3 location>\Saves and <civ3 location>\Saves\Auto to watch list` + "\n")

	// Add Saves and Saves\Auto folder watches
	err = savWatcher.Add(civPath + `\Saves`)
	if err != nil {
		log.Fatal(err)
	}
	err = savWatcher.Add(civPath + `\Saves\Auto`)
	if err != nil {
		log.Fatal(err)
	}

	// Api server
	go server()

	// _, err = lorca.New("Civ Intelligence Agency III", "http://"+addr+":"+port+"/isocss.html", 800, 600)
	ui, err := lorca.New("", "", 1280, 720)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	ui.Load("http://" + addr + ":" + port + "/isocss.html")
	<-ui.Done()
}
