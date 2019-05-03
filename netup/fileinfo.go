package netup

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path"
	"strings"
)

type compareResult int

const (
	fileNotExist = iota + 1
	fileNotMatch
	fileIsMatch
)

const fileInfoFormat = "%s\t%d\t%08x\n"

type fileInfo struct {
	name string
	size uint64
	hash uint32
}

func (info fileInfo) compareWithFile(name string) (compareResult, error) {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return fileNotExist, nil
		}
		return 0, err
	}
	if (uint64)(fi.Size()) != info.size {
		return fileNotMatch, nil
	}
	v, err := calcCRC32(name)
	if err != nil {
		if os.IsNotExist(err) {
			return fileNotExist, nil
		}
		return 0, err
	}
	if v != info.hash {
		return fileNotMatch, nil
	}
	return fileIsMatch, nil
}

func (info fileInfo) dir() string {
	return path.Dir(info.name)
}

func (info fileInfo) dirList() []string {
	return strings.Split(info.dir(), "/")
}

type fileInfoTable map[string]fileInfo

func loadFileInfo(fname string) (fileInfoTable, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := bufio.NewReader(f)
	t := make(fileInfoTable)
	for {
		l, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		fi := fileInfo{}
		if _, err := fmt.Sscanf(l, fileInfoFormat, &fi.name, &fi.size, &fi.hash); err != nil {
			return nil, err
		}
		t[fi.name] = fi
	}
	return t, nil
}

func calcCRC32(name string) (uint32, error) {
	r, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	h := crc32.NewIEEE()
	if _, err := io.Copy(h, r); err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
