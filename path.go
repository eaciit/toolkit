package toolkit

import (
	"os"
	"path/filepath"
)

var (
	PathSeparator string = string(os.PathSeparator)
)

func PathDefault(removeSlash bool) string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if removeSlash == false {
		dir = dir + "/"
	}
	return dir
}
