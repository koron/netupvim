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
	errSourceNotFound    = errors.New("source not found")
	errSourceNotModified = errors.New("source not modified")

	errGithubNoRelease       = errors.New("absence of github release")
	errGithubNoAssets        = errors.New("no matched assets in github release")
	errGithubIncompleteAsset = errors.New("incomplete github asset")
)

type progressFunc func(curr, max int64)

type source interface {
	// download downloads source file to outdir, return its path name.
	// if pivot is not zero, this checks changes of source from pivot.
	download(outdir string, pivot time.Time, f progressFunc) (path string, err error)

	stripCount() int

	// String returns a string to represent source.
	String() string
}

type directSource struct {
	url   string
	strip int
}

var _ source = (*directSource)(nil)

func (ds *directSource) download(d string, p time.Time, f progressFunc) (string, error) {
	return download(ds.url, d, p, f)
}

func (ds *directSource) stripCount() int {
	return ds.strip
}

func (ds *directSource) String() string {
	return fmt.Sprintf("direct: URL=%s", ds.url)
}

type githubSource struct {
	user    string
	project string
	namePat *regexp.Regexp
	strip   int
}

var _ source = (*githubSource)(nil)

func (gs *githubSource) download(d string, p time.Time, f progressFunc) (string, error) {
	a, err := gs.fetchAsset()
	if err != nil {
		return "", err
	}
	if !p.IsZero() && p.After(a.UpdatedAt) {
		return "", errSourceNotModified
	}
	msgPrintln("found newer release on GitHub")
	return download(a.DownloadURL, d, p, f)
}

func (gs *githubSource) stripCount() int {
	return gs.strip
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

func (gs *githubSource) String() string {
	return fmt.Sprintf("GitHub: %s/%s pattern=%s",
		gs.user, gs.project, gs.namePat.String())
}

func downloadFilepath(inURL, outdir string) (string, error) {
	u, err := url.Parse(inURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(outdir, filepath.Base(u.Path)), nil
}

var downloadTimeout = 5 * time.Minute

func downloadAsFile(inURL, outPath string, pivot time.Time, pf progressFunc) error {
	req, err := http.NewRequest("GET", inURL, nil)
	if err != nil {
		return err
	}
	if !pivot.IsZero() {
		t := pivot.UTC().Format(http.TimeFormat)
		req.Header.Set("If-Modified-Since", t)
	}
	logInfo("download URL %s as file %s", inURL, outPath)
	msgPrintf("download %s\n", inURL)
	client := http.Client{Timeout: downloadTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return saveBody(outPath, resp, pf)
	case http.StatusNotModified:
		return errSourceNotModified
	default:
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}
}

// download downloads URL and saves as a file to outdir, return its path name.
// if pivot is not zero, this checks changes of source after pivot.
func download(inURL, outdir string, pivot time.Time, f progressFunc) (string, error) {
	path, err := downloadFilepath(inURL, outdir)
	if err != nil {
		return "", err
	}
	if err := downloadAsFile(inURL, path, pivot, f); err != nil {
		return "", err
	}
	return path, nil
}

func saveBody(outPath string, resp *http.Response, pf progressFunc) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := &progressWriter{w: f, f: pf, m: resp.ContentLength}
	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return nil
}

type progressWriter struct {
	w    io.Writer
	f    progressFunc
	n, m int64
}

func (w *progressWriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.n += int64(n)
	if w.f != nil {
		w.f(w.n, w.m)
	}
	return n, err
}

// sourceSet is set of source.
type sourceSet map[string]map[arch.CPU]source

// Find finds a source for type and CPU.
func (ss sourceSet) Find(sourceType string, cpu arch.CPU) (source, error) {
	m, ok := ss[sourceType]
	if !ok {
		return nil, errSourceNotFound
	}
	s, ok := m[cpu]
	if !ok {
		return nil, errSourceNotFound
	}
	return s, nil
}
