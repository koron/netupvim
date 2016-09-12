package main

import (
	"regexp"

	"github.com/koron/go-arch"
	"github.com/koron/netupvim/netup"
)

var vimSet = map[string]netup.SourcePack{
	"release": {
		arch.X86: &netup.GithubSource{
			Name:    "vim",
			User:    "koron",
			Project: "vim-kaoriya",
			NamePat: regexp.MustCompile(`-win32-.*\.zip$`),
			Strip:   1,
		},
		arch.AMD64: &netup.GithubSource{
			Name:    "vim",
			User:    "koron",
			Project: "vim-kaoriya",
			NamePat: regexp.MustCompile(`-win64-.*\.zip$`),
			Strip:   1,
		},
	},
	"develop": {
		arch.X86: &netup.DirectSource{
			Name:  "vim",
			URL:   "https://files.kaoriya.net/vim/vim-kaoriya-win32-develop.zip",
			Strip: 1,
		},
		arch.AMD64: &netup.DirectSource{
			Name:  "vim",
			URL:   "https://files.kaoriya.net/vim/vim-kaoriya-win64-develop.zip",
			Strip: 1,
		},
	},
	"canary": {
		arch.X86: &netup.DirectSource{
			Name:  "vim",
			URL:   "https://files.kaoriya.net/vim/vim-kaoriya-win32-canary.zip",
			Strip: 1,
		},
		arch.AMD64: &netup.DirectSource{
			Name:  "vim",
			URL:   "https://files.kaoriya.net/vim/vim-kaoriya-win64-canary.zip",
			Strip: 1,
		},
	},
	"vim.org": {
		arch.X86: &netup.GithubSource{
			Name:    "vim",
			User:    "vim",
			Project: "vim-win32-installer",
			NamePat: regexp.MustCompile(`_x86\.zip$`),
			Strip:   2,
		},
		arch.AMD64: &netup.GithubSource{
			Name:    "vim",
			User:    "vim",
			Project: "vim-win32-installer",
			NamePat: regexp.MustCompile(`_x64\.zip$`),
			Strip:   2,
		},
	},
}

var netupPack = netup.SourcePack{
	arch.X86: &netup.GithubSource{
		Name:    "netup",
		User:    "koron",
		Project: "netupvim",
		NamePat: regexp.MustCompile(`^netupvim-.*\.zip$`),
	},
}
