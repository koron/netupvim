package netup

type Package struct {
	Strip int
}

type DirectPackage struct {
}

type GithubPacakge struct {
	User    string
	Project string
}

// Update updates or installs a package into target directory.
func Update(dir string, pkg Package) error {
	// TODO:
	return nil
}
