package config

import (
	_ "embed"
	"encoding/json"
	"github.com/sebastianappelberg/disk/pkg/cache"
)

type FolderSet map[string]bool

type configItem struct {
	Description string   `json:"description"`
	Folders     []string `json:"folders"`
}

var (
	//go:embed clutter_folders.json
	clutterFoldersJSON []byte

	//go:embed unsafe_folders.json
	unsafeFoldersJSON []byte

	ClutterFolders      FolderSet
	UnsafeFolders       FolderSet
	UserExcludedFolders FolderSet
	configCache         *cache.Cache[FolderSet]

	userExcludedFoldersKey = "excludedFolders"
)

func init() {
	setConfig()
}

func setConfig() {
	ClutterFolders = mustGetFolderConfig(clutterFoldersJSON)
	UnsafeFolders = mustGetFolderConfig(unsafeFoldersJSON)
	configCache = cache.NewCache[FolderSet](GetAppDir(), "user_config")
	if folderSet, ok := configCache.Get(userExcludedFoldersKey); ok {
		UserExcludedFolders = folderSet
	} else {
		UserExcludedFolders = make(map[string]bool)
	}
}

func ExcludeFolder(path string) {
	UserExcludedFolders[path] = true
	configCache.Put(userExcludedFoldersKey, UserExcludedFolders)
	configCache.Flush()
}

// Load JSON configuration
func parseFolderConfig(data []byte) (map[string]configItem, error) {
	var config map[string]configItem
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getFolderSet(config map[string]configItem) FolderSet {
	folders := make(FolderSet)
	for _, category := range config {
		for _, folder := range category.Folders {
			folders[folder] = true
		}
	}
	return folders
}

func mustGetFolderConfig(data []byte) FolderSet {
	parsed, err := parseFolderConfig(data)
	if err != nil {
		panic(err)
	}
	return getFolderSet(parsed)
}
