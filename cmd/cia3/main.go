package main

import (
	// "fmt"

	"net"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jcuga/golongpoll"
	"github.com/zserge/lorca"
)

var savWatcher *fsnotify.Watcher
var debounceTimer *time.Timer
var longPoll *golongpoll.LongpollManager
var listener net.Listener
var errorChannel = make(chan error, 20)

const debounceInterval = 300 * time.Millisecond

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

	var lastSav string
	// Read Win registry for Civ3 Conquests path
	civPath, err := findWinCivInstall()
	if err == nil {
		err = loadDefaultBiq(civPath + `\conquests.biq`)
		if err != nil {
			errorChannel <- err
		}
		lastSav, err = getLastSav(civPath)
		if err != nil {
			errorChannel <- err
		} else {
			err = loadNewSav(lastSav)
			if err != nil {
				errorChannel <- err
			}
		}
		// Add Saves and Saves\Auto folder watches
		err = savWatcher.Add(civPath + `\Saves`)
		if err != nil {
			errorChannel <- err
		}
		err = savWatcher.Add(civPath + `\Saves\Auto`)
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

	ui, err := lorca.New("", "", 1280, 720)
	if err == nil {
		defer ui.Close()
		ui.Load(httpUrlString)
		<-ui.Done()
	} else {
		errorChannel <- err
		// fallback to fyne GUI with hyperlink
		fyneUi()
	}
}
