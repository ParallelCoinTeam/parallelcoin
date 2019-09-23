package util

import (
	"os"
	"path/filepath"
	"strings"
)

// EnsureDir checks a file could be written to a path,
// creates the directories as needed
func EnsureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// MinUint32 is a helper function to return the minimum of two uint32s.
// This avoids a math import and the need to cast to floats.
func MinUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// RightJustify takes a string and right justifies it by a width or crops it
func RightJustify(s string, w int) string {
	sw := len(s)
	diff := w - sw
	if diff > 0 {
		s = strings.Repeat(" ", diff) + s
	} else if diff < 0 {
		s = s[:w]
	}
	return s
}
