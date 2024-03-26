package commands

import (
	"github.com/stretchr/testify/assert"
	"hjbdev/pvm/common"
	"io/fs"
	"os"
	"testing"
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

func Test_sortVersions_sortsVersionsDescending(t *testing.T) {
	input := []versionMeta{
		{
			number: common.Version{
				Major: "7",
				Minor: "4",
				Patch: "1",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "7",
				Minor: "4",
				Patch: "2",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "8",
				Minor: "3",
				Patch: "0",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "8",
				Minor: "3",
				Patch: "1",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "8",
				Minor: "2",
				Patch: "0",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "8",
				Minor: "2",
				Patch: "5",
			},
			folder: fakeDirEntry{},
		},
	}

	output := sortVersions(input)

	assert.Equal(t, []versionMeta{
		{
			number: common.Version{
				Major: "8",
				Minor: "3",
				Patch: "1",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "8",
				Minor: "3",
				Patch: "0",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "8",
				Minor: "2",
				Patch: "5",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "8",
				Minor: "2",
				Patch: "0",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "7",
				Minor: "4",
				Patch: "2",
			},
			folder: fakeDirEntry{},
		},
		{
			number: common.Version{
				Major: "7",
				Minor: "4",
				Patch: "1",
			},
			folder: fakeDirEntry{},
		},
	}, output)
}
