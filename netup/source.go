package netup

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
	errSourceNotModified = errors.New("source not modified")

	errGithubNoRelease       = errors.New("absence of github release")
	errGithubNoAssets        = errors.New("no matched assets in github release")
	errGithubIncompleteAsset = errors.New("incomplete github asset")
)

type progressFunc func(curr, max int64)

// Source describes source of update.
type Source interface {
	// download downloads source file to outdir, return its path name.
	// if pivot is not zero, this checks changes of source from pivot.
	download(outdir string, pivot time.Time, f progressFunc) (path string, err error)

	stripCount() int

	name() string

	// String returns a string to represent source.
	String() string
}

// DirectSource represents direct ZIP source.
type DirectSource struct {
	Name  string
	URL   string
	Strip int
}

var _ Source = (*DirectSource)(nil)

func (ds *DirectSource) download(d string, p time.Time, f progressFunc) (string, error) {
	return download(ds.URL, d, p, f)
}

func (ds *DirectSource) stripCount() int {
	return ds.Strip
}

func (ds *DirectSource) name() string {
	return ds.Name
}

func (ds *DirectSource) String() string {
	return fmt.Sprintf("direct: URL=%s", ds.URL)
}

// GithubSource represents project source on GitHub.
type GithubSource struct {
	Name    string
	User    string
	Project string
	NamePat *regexp.Regexp
	Strip   int
}

var _ Source = (*GithubSource)(nil)

func (gs *GithubSource) download(d string, p time.Time, f progressFunc) (string, error) {
	a, err := gs.fetchAsset(p)
	if err != nil {
		return "", err
	}
	if !p.IsZero() && p.After(a.UpdatedAt) {
		return "", errSourceNotModified
	}
	msgPrintln("found newer release on GitHub")
	return download(a.DownloadURL, d, p, f)
}

func (gs *GithubSource) stripCount() int {
	return gs.Strip
}

func (gs *GithubSource) name() string {
	return gs.Name
}

func (gs *GithubSource) fetchAsset(pivot time.Time) (*github.Asset, error) {
	r, err := github.LatestIfModifiedSince(gs.User, gs.Project, pivot)
	if err != nil {
		if err == github.ErrNotModified {
			err = errSourceNotModified
		}
		return nil, err
	}
	if r.Draft || r.PreRelease {
		return nil, errGithubNoRelease
	}
	var t *github.Asset
	for _, a := range r.Assets {
		if gs.NamePat.MatchString(a.Name) {
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

func (gs *GithubSource) String() string {
	return fmt.Sprintf("GitHub: %s/%s pattern=%s",
		gs.User, gs.Project, gs.NamePat.String())
}

func downloadFilepath(inURL, outdir string) (string, error) {
	u, err := url.Parse(inURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(outdir, filepath.Base(u.Path)), nil
}

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
	client := http.Client{Timeout: DownloadTimeout}
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

// SourcePack is the map arch.CPU to source.
type SourcePack map[arch.CPU]Source
