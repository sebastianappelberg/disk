package games

import (
	"errors"
	"golang.org/x/sys/windows/registry"
	"os"
	"runtime"
)

func getSteamPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "~/Library/Application Support/Steam"
	case "windows":
		path, err := getWindowsSteamPath()
		if err != nil {
			panic(err)
		}
		return path
	case "linux":
		path, err := getLinuxSteamPath()
		if err != nil {
			panic(err)
		}
		return path
	default:
		panic("No steam folder found")
	}
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
