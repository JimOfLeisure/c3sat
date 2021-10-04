//go:build !windows

package luaciv3

// TODO: Insert other logic, perhaps take parameter or read env key
func findCivInstall() (string, error) {
	return "", nil
}
