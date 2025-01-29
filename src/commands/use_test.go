package commands

import (
	"io/fs"
	"os"
)

// mock os.DirEntry for test
type fakeDirEntry struct{}

func (f fakeDirEntry) Name() string {
	return ""
}

func (f fakeDirEntry) IsDir() bool {
	return true
}

func (f fakeDirEntry) Type() os.FileMode {
	return fs.ModeDir
}

func (f fakeDirEntry) Info() (os.FileInfo, error) {
	return nil, nil
}
