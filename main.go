package main

import (
	"regexp"

	"github.com/koron/go-arch"
)

var sources = sourceSet{
	"release": {
		arch.X86: &githubSource{
			user:    "koron",
			project: "vim-kaoriya",
			namePat: regexp.MustCompile(`-win32-.*\.zip$`),
			strip:   1,
		},
		arch.AMD64: &githubSource{
			user:    "koron",
			project: "vim-kaoriya",
			namePat: regexp.MustCompile(`-win64-.*\.zip$`),
			strip:   1,
		},
	},
	"develop": {
		arch.X86: &directSource{
			url:   "http://files.kaoriya.net/vim/vim74-kaoriya-win32.zip",
			strip: 1,
		},
		arch.AMD64: &directSource{
			url:   "http://files.kaoriya.net/vim/vim74-kaoriya-win64.zip",
			strip: 1,
		},
	},
	"canary": {
		arch.X86: &directSource{
			url:   "http://files.kaoriya.net/vim/vim74-kaoriya-win32-test.zip",
			strip: 1,
		},
		arch.AMD64: &directSource{
			url:   "http://files.kaoriya.net/vim/vim74-kaoriya-win64-test.zip",
			strip: 1,
		},
	},
	"vim.org": {
		arch.X86: &githubSource{
			user:    "vim",
			project: "vim-win32-installer",
			namePat: regexp.MustCompile(`_x86\.zip$`),
			strip:   2,
		},
		arch.AMD64: &githubSource{
			user:    "vim",
			project: "vim-win32-installer",
			namePat: regexp.MustCompile(`_x64\.zip$`),
			strip:   2,
		},
	},
}

func main() {
	err := run("netupvim.ini", sources)
	if err != nil {
		logFatal(err)
	}
}
