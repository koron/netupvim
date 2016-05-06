package main

import (
	"regexp"

	"github.com/koron/go-arch"
	"github.com/koron/netupvim/netup"
)

var sources = netup.SourceSet{
	"release": {
		arch.X86: &netup.GithubSource{
			User:    "koron",
			Project: "vim-kaoriya",
			NamePat: regexp.MustCompile(`-win32-.*\.zip$`),
			Strip:   1,
		},
		arch.AMD64: &netup.GithubSource{
			User:    "koron",
			Project: "vim-kaoriya",
			NamePat: regexp.MustCompile(`-win64-.*\.zip$`),
			Strip:   1,
		},
	},
	"develop": {
		arch.X86: &netup.DirectSource{
			URL:   "http://files.kaoriya.net/vim/vim74-kaoriya-win32.zip",
			Strip: 1,
		},
		arch.AMD64: &netup.DirectSource{
			URL:   "http://files.kaoriya.net/vim/vim74-kaoriya-win64.zip",
			Strip: 1,
		},
	},
	"canary": {
		arch.X86: &netup.DirectSource{
			URL:   "http://files.kaoriya.net/vim/vim74-kaoriya-win32-test.zip",
			Strip: 1,
		},
		arch.AMD64: &netup.DirectSource{
			URL:   "http://files.kaoriya.net/vim/vim74-kaoriya-win64-test.zip",
			Strip: 1,
		},
	},
	"vim.org": {
		arch.X86: &netup.GithubSource{
			User:    "vim",
			Project: "vim-win32-installer",
			NamePat: regexp.MustCompile(`_x86\.zip$`),
			Strip:   2,
		},
		arch.AMD64: &netup.GithubSource{
			User:    "vim",
			Project: "vim-win32-installer",
			NamePat: regexp.MustCompile(`_x64\.zip$`),
			Strip:   2,
		},
	},
}

func main() {
	err := netup.Run("netupvim.ini", sources)
	if err != nil {
		netup.LogFatal(err)
	}
}
