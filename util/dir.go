package util

import (
	"os"
	"path/filepath"
)

// TraversalDirByExt Traversal dir file
func TraversalDirByExt(dir, ext string, fc func(fp string)) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			fn := info.Name()
			if filepath.Ext(fn) == ext {
				fc(path)
			}
		}
		return nil
	})
}

// MustBeDir check dir exist, create if not exist
func MustBeDir(dir string) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0777); err != nil {
				panic("dir not exist and create failed:" + dir + ", err: " + err.Error())
			}
		} else {
			panic("stat file: " + dir + " failed")
		}
		return
	}

	if !info.IsDir() {
		panic("a file exist with the same name: " + dir)
	}
}
