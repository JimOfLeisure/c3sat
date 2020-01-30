package main

import (
	"net"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jcuga/golongpoll"
)

// build with `-ldflags "-X main.appVersion=myVersionName"` to set version at compile time
var appVersion = "0.4.1-devbuild"

var savWatcher *fsnotify.Watcher
var debounceTimer *time.Timer
var longPoll *golongpoll.LongpollManager
var listener net.Listener
var errorChannel = make(chan error, 20)
var watchList = new(watchListType)

const debounceInterval = 300 * time.Millisecond

func loadStartupFiles(civPath string) {
	var err error
	err = loadDefaultBiq(civPath + `\conquests.biq`)
	if err != nil {
		errorChannel <- err
	}
	lastSav, err := getLastSav(civPath)
	if err != nil {
		errorChannel <- err
	} else {
		err = loadNewSav(lastSav)
		if err != nil {
			errorChannel <- err
		}
	}
}

func main() {
	var err error
	savWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		errorChannel <- err
	}
	defer savWatcher.Close()

	// Set up file event handler
	debounceTimer = time.NewTimer(debounceInterval)
	go watchSavs()

	// Initialize long poll manager
	longPoll, err = golongpoll.StartLongpoll(golongpoll.Options{})
	if err != nil {
		errorChannel <- err
	}
	defer longPoll.Shutdown()

	// Read Win registry for Civ3 Conquests path
	civPath, err := findWinCivInstall()
	if err == nil {
		loadStartupFiles(civPath)
		// Add Saves and Saves\Auto folder watches
		err = watchList.addWatch(civPath + `\Saves`)
		if err != nil {
			errorChannel <- err
		}
		err = watchList.addWatch(civPath + `\Saves\Auto`)
		if err != nil {
			errorChannel <- err
		}

	} else {
		errorChannel <- err
	}

	for i := 0; i < len(httpPortTry); i++ {
		listener, err = net.Listen("tcp", httpPortTry[i])
		if err == nil {
			break
		}
	}
	if listener == nil {
		panic(err)
	}

	httpPort = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	httpUrlString = "http://" + addr + ":" + httpPort + "/"

	// Api server
	go server()

	go func() {
		for {
			select {
			case err := <-errorChannel:
				longPoll.Publish("exeError", err.Error())
			}
		}
	}()

	// fyne GUI with hyperlink
	fyneUi()
}
