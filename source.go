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
	"github.com/koron/go-github"
)

var (
	errUnknownSource     = errors.New("unknown source")
	errSourceNotFound    = errors.New("source not found")
	errSourceNotModified = errors.New("source not modified")

	errGithubNoRelease       = errors.New("absence of github release")
	errGithubNoAssets        = errors.New("no matched assets in github release")
	errGithubIncompleteAsset = errors.New("incomplete github asset")
)

type sourceType int

const (
	releaseSource sourceType = iota
	developSource
	canarySource
)

func toSourceType(s string) (sourceType, error) {
	switch s {
	case "release":
		return releaseSource, nil
	case "develop":
		return developSource, nil
	case "canary":
		return canarySource, nil
	}
	return 0, errUnknownSource
}

type source interface {
	// download downloads source file to outdir, return its path name.
	// if pivot is not zero, this checks changes of source from pivot.
	download(outdir string, pivot time.Time) (path string, err error)
}

type directSource struct {
	url string
}

var _ source = (*directSource)(nil)

func (ds *directSource) download(d string, p time.Time) (string, error) {
	return download(d, ds.url, p)
}

type githubSource struct {
	user    string
	project string
	namePat *regexp.Regexp
}

var _ source = (*githubSource)(nil)

func (gs *githubSource) download(d string, p time.Time) (string, error) {
	a, err := gs.fetchAsset()
	if err != nil {
		return "", err
	}
	if !p.IsZero() && p.After(a.UpdatedAt) {
		return "", errSourceNotModified
	}
	return download(a.DownloadURL, d, p)
}

func (gs *githubSource) fetchAsset() (*github.Asset, error) {
	r, err := github.Latest(gs.user, gs.project)
	if err != nil {
		return nil, err
	}
	if r.Draft || r.PreRelease {
		return nil, errGithubNoRelease
	}
	var t *github.Asset
	for _, a := range r.Assets {
		if gs.namePat.MatchString(a.Name) {
			t = &a
			break
		}
	}
	if t == nil {
		return nil, errGithubNoAssets
	}
	if t.State != "uploaded" {
		return nil, errGithubIncompleteAsset
	}
	return t, nil
}

var sources = map[sourceType]map[arch.CPU]source{
	releaseSource: {
		arch.X86: &githubSource{
			user:    "koron",
			project: "vim-kaoriya",
			namePat: regexp.MustCompile(`-win32-.*\.zip$`),
		},
		arch.AMD64: &githubSource{
			user:    "koron",
			project: "vim-kaoriya",
			namePat: regexp.MustCompile(`-win64-.*\.zip$`),
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

func downloadFilepath(inURL, outdir string) (string, error) {
	u, err := url.Parse(inURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(outdir, filepath.Base(u.Path)), nil
}

func downloadAsFile(inURL, outPath string, pivot time.Time) error {
	req, err := http.NewRequest("GET", inURL, nil)
	if err != nil {
		return err
	}
	if !pivot.IsZero() {
		t := pivot.UTC().Format(http.TimeFormat)
		req.Header.Set("If-Modified-Since", t)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, resp.Body); err != nil {
			return err
		}
	case http.StatusNotModified:
		return errSourceNotModified
	default:
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}
	return nil
}

// download downloads URL and saves as a file to outdir, return its path name.
// if pivot is not zero, this checks changes of source after pivot.
func download(inURL, outdir string, pivot time.Time) (string, error) {
	path, err := downloadFilepath(inURL, outdir)
	if err != nil {
		return "", err
	}
	if err := downloadAsFile(inURL, path, pivot); err != nil {
		return "", err
	}
	return path, nil
}
