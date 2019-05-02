package rotate

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var Logger *log.Logger = nil

func logf(format string, v ...interface{}) {
	if Logger == nil {
		return
	}
	Logger.Printf(format, v...)
}

// Rotate evacuate a file with rotations.
func Rotate(name string, max int) error {
	last := rotateName(name, max)
	err := os.Remove(last)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for i := max - 1; i >= 0; i-- {
		curr := rotateName(name, i)
		err := os.Rename(curr, last)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		last = curr
	}
	return nil
}

func rotateName(name string, index int) string {
	if index == 0 {
		return name
	}
	base, ext := splitExt(name)
	return base + "." + strconv.Itoa(index) + ext
}

func splitExt(name string) (base, ext string) {
	ext = path.Ext(name)
	base = name[:len(name)-len(ext)]
	return base, ext
}

var targets = map[string]struct{}{
	".exe": struct{}{},
	".dll": struct{}{},
}

// IsTarget check whether the file is rotate target or not.
func IsTarget(name string) bool {
	ext := strings.ToLower(path.Ext(name))
	_, ok := targets[ext]
	return ok
}

// Sweep deletes rotated files.
func Sweep(name string, max int) error {
	first := true
	for i := 1; i <= max; i++ {
		n := rotateName(name, i)
		fi, err := os.Stat(n)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if fi.IsDir() {
			continue
		}
		err = os.Remove(n)
		if err != nil {
			return err
		}
		if first {
			logf("sweep for %s", name)
			first = false
		}
		logf("- deleted %s", n)
	}
	return nil
}
