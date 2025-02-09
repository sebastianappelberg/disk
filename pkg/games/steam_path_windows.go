//go:build windows

package games

import (
	"golang.org/x/sys/windows/registry"
)

func getSteamPath() string {
	path, err := getWindowsSteamPath()
	if err != nil {
		panic(err)
	}
	return path
}

func getWindowsSteamPath() (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		// Try fallback for 32-bit Windows
		key, err = registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
		if err != nil {
			return "", err
		}
	}
	defer key.Close()

	steamPath, _, err := key.GetStringValue("InstallPath")
	if err != nil {
		return "", err
	}
	return steamPath, nil
}
