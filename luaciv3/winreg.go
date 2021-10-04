//go:build windows

package luaciv3

import (
	"golang.org/x/sys/windows/registry"
)

// really should be const, but can't have literal string arrays as const
var civInstallPathKeyTry = []string{
	`SOFTWARE\WOW6432Node\Infogrames\Conquests`,
	`SOFTWARE\Infogrames\Conquests`,
}

func findCivInstall() (string, error) {
	var k registry.Key
	var err error
	for i := 0; i < len(civInstallPathKeyTry); i++ {
		k, err = registry.OpenKey(registry.LOCAL_MACHINE, civInstallPathKeyTry[i], registry.QUERY_VALUE)
		if err != nil {
			if i >= len(civInstallPathKeyTry)-1 {
				return "", err
			}
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
