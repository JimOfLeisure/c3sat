package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/sys/windows/registry"
)

const civInstallPathKey = `SOFTWARE\WOW6432Node\Infogrames\Conquests`

// assume go prog is 64-bit
func findWinCivInstall() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, civInstallPathKey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
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
	_ = pathEnd
	for i := pathStart; i < (len(ini)); i++ {
		if ini[i] == '\r' || ini[i] == '\n' {
			pathEnd = i - 1
			break
		}
	}
	if pathEnd <= pathStart {
		return "", fmt.Errorf("Failed to find Latest Save in conquests.ini")
	}
	return string(ini[pathStart:pathEnd]) + ".SAV", nil
}
