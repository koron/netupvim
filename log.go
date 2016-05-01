package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger = log.New(ioutil.Discard, "", 0)

// logInfo records a message to logger file.
func logInfo(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Println(s)
}

// logWarn records a message to UI and logger file.
func logWarn(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	msgPrintln(s)
	logger.Println(s)
}

// logFatal records a message to UI and logger file then os.Exit(1)
func logFatal(err error) {
	msgPrintln(err)
	logger.Println(err)
	os.Exit(1)
}

func logLoadRecipeFailed(err error) {
	if os.IsExist(err) {
		logWarn("failed to load recipe, try to extract all files: %s", err)
	}
}

func logSaveRecipeFailed(err error) {
	logWarn("failed to save recipe: %s", err)
}

func logCompareFileFailed(err error, name string) {
	logWarn("failed to compare file %q: %s", name, err)
}

func logCleanArchiveFailed(err error) {
	logWarn("failed to remove downloaded archive: %s", err)
}

func logCleanLogFailed(err error) {
	logWarn("failed to remove old log file: %s", err)
}

const logLayout = "20060102T150405Z0700.log"

func logFiles(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var logs []os.FileInfo
	for _, fi := range files {
		_, err := time.Parse(logLayout, fi.Name())
		if err != nil {
			continue
		}
		logs = append(logs, fi)
	}
	return logs, nil
}

func logSetup(dir string, count int) {
	// remove old log files.
	logs, err := logFiles(dir)
	if len(logs) >= count {
		for _, fi := range logs[:len(logs)-count+1] {
			err := os.Remove(filepath.Join(dir, fi.Name()))
			if err != nil {
				logCleanLogFailed(err)
			}
		}
	}
	// create a logger with new log file.
	name := filepath.Join(dir, time.Now().Format(logLayout))
	f, err := os.Create(name)
	if err != nil {
		logFatal(err)
	}
	logger = log.New(f, "", log.LstdFlags)
}
