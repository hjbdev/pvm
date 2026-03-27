package common

import (
	"os"
	"path/filepath"
)

type PVMPaths struct {
	Home               string
	Root               string
	VersionsDir        string
	BinDir             string
	CurrentVersionFile string
}

func NewPVMPaths() (PVMPaths, error) {
	// Allow tests to override the home directory by setting HOME.
	// On Windows, os.UserHomeDir uses USERPROFILE; prefer HOME if present to
	// make tests platform-independent.
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return PVMPaths{}, err
		}
	}

	root := filepath.Join(homeDir, ".pvm")
	return PVMPaths{
		Home:               homeDir,
		Root:               root,
		VersionsDir:        filepath.Join(root, "versions"),
		BinDir:             filepath.Join(root, "bin"),
		CurrentVersionFile: filepath.Join(root, "version"),
	}, nil
}

func (p PVMPaths) VersionDir(name string) string {
	return filepath.Join(p.VersionsDir, name)
}
