package config

import (
	_ "embed"
	"github.com/sebastianappelberg/disk/pkg/util"
	"os"
)

var (
	appDir string
)

func init() {
	setAppDir()
}

func setAppDir() {
	dirFromEnv := os.Getenv("DISK_DIR")
	if dirFromEnv != "" {
		appDir = dirFromEnv
		return
	}
	appDirName := ".disk"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		appDir = appDirName
		return
	}
	appDir = util.SimpleJoin(homeDir, appDirName)
}

func GetAppDir() string {
	return appDir
}
