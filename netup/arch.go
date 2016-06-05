package netup

import (
	"path/filepath"

	"github.com/koron/go-arch"
)

// Arch determines architecture.
type Arch struct {
	// Name for architecture, like "X86", "AMD64"
	Name string

	// Hint is a file be used to guess architecture.
	Hint string
}

func (a *Arch) detectCPU(dir string) (arch.CPU, error) {
	cpu := arch.ParseCPU(a.Name)
	if cpu != 0 {
		return cpu, nil
	}
	return arch.Exe(filepath.Join(dir, a.Hint))
}
