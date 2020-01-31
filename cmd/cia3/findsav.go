package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/myjimnelson/c3sat/queryciv3"

	"golang.org/x/sys/windows/registry"
)

// really should be const, but can't have literal string arrays as const
var civInstallPathKeyTry = []string{
	`SOFTWARE\WOW6432Node\Infogrames\Conquests`,
	`SOFTWARE\Infogrames\Conquests`,
}

func findWinCivInstall() (string, error) {
	var k registry.Key
	var err error
	for i := 0; i < len(civInstallPathKeyTry); i++ {
		k, err = registry.OpenKey(registry.LOCAL_MACHINE, civInstallPathKeyTry[i], registry.QUERY_VALUE)
		if err != nil {
			return "", err
		} else {
			break
		}
	}
	defer k.Close()
	s, _, err := k.GetStringValue("install_path")
	if err != nil {
		return "", err
	}
	return s, nil
}

// Look in conquests.ini in the Civ3 Conquests install path for the "Lastest Save=" value
func getLastSav(path string) (string, error) {
	const key = `Latest Save=`
	file, err := os.Open(path + `\conquests.ini`)
	if err != nil {
		return "", err
	}
	ini, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var pathStart, pathEnd int
	for i := 0; i < (len(ini) - len(key)); i++ {
		if ini[i] == key[0] {
			if string(ini[i:i+len(key)]) == key {
				pathStart = i + len(key)
				break
			}
		}
	}
	for i := pathStart; i < (len(ini)); i++ {
		if ini[i] == '\r' || ini[i] == '\n' {
			pathEnd = i
			break
		}
	}
	if pathEnd <= pathStart {
		return "", fmt.Errorf("Failed to find Latest Save in conquests.ini")
	}
	//  Assuming conquests.ini file is not UTF-8
	s, err := queryciv3.CivString(ini[pathStart:pathEnd])
	// My .ini has a space after the filename, so trimming leading/trailing whitespace
	return strings.TrimSpace(s) + ".SAV", err
}
