//go:build darwin

package games

func getSteamPath() string {
	return "~/Library/Application Support/Steam"
}
