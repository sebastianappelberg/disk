package games

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/vdf"
	"github.com/sebastianappelberg/disk/pkg/util"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type SteamGame struct {
	AppId        string        // AppId is the steamID of the game.
	Name         string        // Name is the name of the game. Needed mainly for display purposes.
	Path         string        // Path is the absolute path to the installation folder.
	ManifestPath string        // ManifestPath is the absolute path to the game's manifest file. It needs to be removed when removing the game for steam to know that the game has been uninstalled.
	Playtime     time.Duration // Playtime is in minutes.
	Size         int64         // Size in bytes.
	LastPlayed   time.Time     // LastPlayed is the date and time of the last time the game was played.
}

func (s SteamGame) GetPaths() []string {
	return []string{s.Path, s.ManifestPath}
}

func getSteamGames() ([]SteamGame, error) {
	games, err := findGames()
	if err != nil {
		return nil, err
	}
	steamUserIds, err := findSteamUserIds()
	if err != nil {
		return nil, err
	}
	configs := getAppConfigs(steamUserIds)
	for i, g := range games {
		appConf := configs[g.AppId]
		games[i].LastPlayed = time.Unix(parseInt(appConf.LastPlayed), 0)
		games[i].Playtime = time.Duration(parseInt(appConf.Playtime)) * time.Minute
	}
	// If there are duplicate paths (which happened in the case of Half-Life/Counter-Strike,
	// prioritize the one with the highest playtime.
	seen := make(map[string]SteamGame)
	for i, g := range games {
		game, ok := seen[g.Path]
		if !ok {
			seen[g.Path] = g
		} else if game.Playtime > g.Playtime {
			games = slices.Delete(games, i, i+1)
		}
	}
	return games, nil
}

type appConfig struct {
	LastPlayed string // LastPlayed is Unix timestamp.
	Playtime   string // Playtime is in minutes.
}

func getAppConfigs(userIds []string) map[string]appConfig {
	configs := make(map[string]appConfig)
	for _, userId := range userIds {
		config, err := getAppConfig(userId)
		if err != nil {
			continue
		}
		for k, v := range config {
			configs[k] = v
		}
	}
	return configs
}

func getAppConfig(userId string) (map[string]appConfig, error) {
	type localConfig struct {
		UserLocalConfigStore struct {
			Software struct {
				Valve struct {
					Steam struct {
						Apps map[string]appConfig `json:"Apps"`
					} `json:"Steam"`
				} `json:"Valve"`
			} `json:"Software"`
		} `json:"UserLocalConfigStore"`
	}

	var config localConfig
	err := readVdfFile(filepath.Join(getSteamPath(), "userdata", userId, "config", "localconfig.vdf"), &config)
	if err != nil {
		return nil, err
	}
	return config.UserLocalConfigStore.Software.Valve.Steam.Apps, nil
}

func findGames() ([]SteamGame, error) {
	folders, err := findGameInstallationFolders()
	if err != nil {
		return nil, err
	}
	var games []SteamGame
	for _, folder := range folders {
		// Read all app manifests.
		appsFolder := util.SimpleJoin(folder, "steamapps")
		dir, err := os.ReadDir(appsFolder)
		if err != nil {
			return nil, err
		}
		for _, file := range dir {
			if strings.HasSuffix(file.Name(), ".acf") {
				manifestPath := util.SimpleJoin(appsFolder, file.Name())
				g, err := readGameConfig(manifestPath)
				if err != nil {
					return nil, err
				}
				g.ManifestPath = manifestPath
				games = append(games, g)
			}
		}
	}
	return games, nil
}

func readGameConfig(path string) (SteamGame, error) {
	type appState struct {
		AppId      string `json:"appid"`      // AppId is the steamID of the game, which is needed when looking up the LastPlayed attribute in localconfig.vdf.
		Name       string `json:"name"`       // Name is the name of the game.
		InstallDir string `json:"installdir"` // InstallDir is the name of the folder the game is stored in, not the entire path.
		SizeOnDisk string `json:"SizeOnDisk"` // SizeOnDisk is the size of the InstallDir in bytes.
	}
	type gameConfigWrapper struct {
		AppState appState `json:"AppState"`
	}

	var g gameConfigWrapper
	err := readVdfFile(path, &g)
	if err != nil {
		return SteamGame{}, err
	}
	result := SteamGame{
		AppId: g.AppState.AppId,
		Path:  filepath.Join(filepath.Dir(path), "common", g.AppState.InstallDir),
		Size:  parseInt(g.AppState.SizeOnDisk),
		Name:  g.AppState.Name,
	}
	return result, nil
}

func findGameInstallationFolders() ([]string, error) {
	type libraryConfig struct {
		Path      string `json:"path"`      // Path is the root of the library. I.e. that's where you'll find your games installed.
		TotalSize string `json:"totalsize"` // TotalSize the size of all the games in the library in bytes. "0" means the library is empty.
	}
	type libraryFoldersWrapper struct {
		LibraryFolders map[string]libraryConfig `json:"libraryfolders"`
	}

	libraryFoldersPath := filepath.Join(getSteamPath(), "steamapps", "libraryfolders.vdf")
	var config libraryFoldersWrapper
	err := readVdfFile(libraryFoldersPath, &config)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w\n", libraryFoldersPath, err)
	}
	var paths []string
	for _, conf := range config.LibraryFolders {
		if conf.TotalSize != "0" {
			paths = append(paths, conf.Path)
		}
	}
	return paths, nil
}

func findSteamUserIds() ([]string, error) {
	type loginUsersWrapper struct {
		Users map[string]interface{} `json:"users"`
	}

	configPath := filepath.Join(getSteamPath(), "config", "loginusers.vdf")
	var config loginUsersWrapper
	err := readVdfFile(configPath, &config)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w\n", configPath, err)
	}
	var result []string
	for id := range config.Users {
		result = append(result, id, toggleSteamIDFormat(id))
	}
	return result, nil
}

func toggleSteamIDFormat(id string) string {
	conversionBase := int64(76561197960265728)
	idInt := parseInt(id)
	if len(id) <= 10 {
		// The id is in SteamID3 format. Convert it to SteamID64 format.
		return strconv.FormatInt(idInt+conversionBase, 10)
	}
	// The id is in SteamID64 format. Convert it to SteamID3 format.
	return strconv.FormatInt(idInt-conversionBase, 10)
}

func parseInt(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func readVdfFile(path string, val any) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	parser := vdf.NewParser(file)
	result, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}
	buffer := &bytes.Buffer{}
	err = json.NewEncoder(buffer).Encode(result)
	if err != nil {
		return fmt.Errorf("error encoding buffer: %w", err)
	}
	err = json.NewDecoder(buffer).Decode(val)
	if err != nil {
		return fmt.Errorf("error decoding buffer: %w", err)
	}
	return nil
}
