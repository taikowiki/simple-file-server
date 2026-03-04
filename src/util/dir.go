package util

import (
	"os"
	"path/filepath"
)

func MainDir() string {
	var mainDir string
	if *Flag.UseCwd {
		cwd, err := os.Getwd()
		if err != nil {
			return ""
		}

		mainDir = cwd
	} else {
		execPath, err := os.Executable()
		if err != nil {
			return ""
		}

		execDir := filepath.Dir(execPath)
		mainDir = execDir
	}
	return mainDir
}

func FileDir() string {
	return filepath.Join(MainDir(), "files")
}
