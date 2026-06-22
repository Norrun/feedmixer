//go:build !(aix || dragonfly || freebsd || (js && wasm) || wasip1 || linux || netbsd || openbsd || solaris)

package data

import (
	"fmt"
	"os"
	"path"
)

func systemPath() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("Error finding system path: %v", err)
	}
	dir := path.Join(home, systemSubDir)
	return dir, nil

}
