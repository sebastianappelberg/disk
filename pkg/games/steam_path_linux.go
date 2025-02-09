//go:build linux

package games

import (
	"errors"
	"os"
)

func getSteamPath() string {
	path, err := getLinuxSteamPath()
	if err != nil {
		panic(err)
	}
	return path
}

func getLinuxSteamPath() (string, error) {
	paths := []string{
		os.ExpandEnv("$HOME/.steam/steam"),
		os.ExpandEnv("$HOME/.local/share/Steam"),
		"/usr/lib/steam",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", errors.New("could not find steam installation path")
}
