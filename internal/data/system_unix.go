//go:build aix || dragonfly || freebsd || (js && wasm) || wasip1 || linux || netbsd || openbsd || solaris

package data

import (
	"os"
	"path"
)

func systemPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := path.Join(home, "/.local/share/", systemSubDir)

	return dir, nil

}
