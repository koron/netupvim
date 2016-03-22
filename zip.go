package main

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func extractZip(zipName, dir string, stripCount int, prev fileInfoTable) (fileInfoTable, error) {
	// extract zip file.
	zr, err := zip.OpenReader(zipName)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	curr := make(fileInfoTable)
	for _, zf := range zr.File {
		if zf.Mode().IsDir() {
			continue
		}
		zfName := stripPath(zf.Name, stripCount)
		curr[zfName] = fileInfo{
			name: zfName,
			size: zf.UncompressedSize64,
			hash: zf.CRC32,
		}
		outName := filepath.Join(dir, zfName)
		// evacuation and optimization.
		if p, ok := prev[zfName]; ok {
			r, err := p.compareWithFile(outName)
			if err != nil {
				// TODO: log an error.
				continue
			}
			switch r {
			case fileNotMatch:
				outName = evacuateName(outName)
			case fileIsMatch:
				// skip un-changed files.
				if p.hash == zf.CRC32 {
					continue
				}
			}
		}
		// rotation.
		ext := strings.ToLower(path.Ext(zfName))
		if ext == ".exe" || ext == ".dll" {
			if err := rotateFiles(outName, 5); err != nil {
				return nil, err
			}
		}
		if err := extractZipFile(zf, outName); err != nil {
			return nil, err
		}
	}
	return curr, nil
}

func extractZipFile(zf *zip.File, name string) error {
	if err := os.MkdirAll(filepath.Dir(name), 0777); err != nil {
		return err
	}
	r, err := zf.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

func stripPath(name string, count int) string {
	s := strings.Split(name, "/")
	return path.Join(s[count:]...)
}

func evacuateName(name string) string {
	base, ext := splitExt(name)
	return base + ".orig" + ext
}

func rotateFiles(name string, max int) error {
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
