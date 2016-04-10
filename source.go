package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/koron/go-arch"
)

type sourceType int

const (
	releaseSource sourceType = iota + 1
	developSource
	canarySource
)

var (
	errSourceNotFound    = errors.New("source not found")
	errSourceNotModified = errors.New("source not modified")
)

type source interface {
	download(outdir string, pivot time.Time) (path string, err error)
}

type directSource struct {
	url string
}

var _ source = (*directSource)(nil)

func (ds *directSource) download(d string, p time.Time) (string, error) {
	out, err := ds.outPath(d)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", ds.url, nil)
	if err != nil {
		return "", err
	}
	if !p.IsZero() {
		t := p.UTC().Format(http.TimeFormat)
		req.Header.Set("If-Modified-Since", t)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		f, err := os.Create(out)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := io.Copy(f, resp.Body); err != nil {
			return "", err
		}

	case http.StatusNotModified:
		return "", errSourceNotModified

	default:
		return "", fmt.Errorf("unexpected response: %s", resp.Status)
	}
	return out, nil
}

func (ds *directSource) outPath(d string) (string, error) {
	u, err := url.Parse(ds.url)
	if err != nil {
		return "", err
	}
	return filepath.Join(d, filepath.Base(u.Path)), nil
}

type githubSource struct {
	user        string
	project     string
	namePattern *regexp.Regexp
}

var _ source = (*githubSource)(nil)

func (gs *githubSource) download(d string, a time.Time) (string, error) {
	// TODO:
	return "", nil
}

var sources = map[sourceType]map[arch.CPU]source{
	releaseSource: {
		arch.X86: &githubSource{
			user:        "koron",
			project:     "vim-kaoriya",
			namePattern: regexp.MustCompile(`-win32-.*\.zip$`),
		},
		arch.AMD64: &githubSource{
			user:        "koron",
			project:     "vim-kaoriya",
			namePattern: regexp.MustCompile(`-win64-.*\.zip$`),
		},
	},
	developSource: {
		arch.X86: &directSource{
			url: "http://files.kaoriya.net/vim/vim74-kaoriya-win32.zip",
		},
		arch.AMD64: &directSource{
			url: "http://files.kaoriya.net/vim/vim74-kaoriya-win64.zip",
		},
	},
	canarySource: {
		arch.X86: &directSource{
			url: "http://files.kaoriya.net/vim/vim74-kaoriya-win32-test.zip",
		},
		arch.AMD64: &directSource{
			url: "http://files.kaoriya.net/vim/vim74-kaoriya-win64-test.zip",
		},
	},
}

func determineSource(st sourceType, cpu arch.CPU) (source, error) {
	m, ok := sources[st]
	if !ok {
		return nil, errSourceNotFound
	}
	s, ok := m[cpu]
	if !ok {
		return nil, errSourceNotFound
	}
	return s, nil
}
