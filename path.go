package dir

// Path returns the full path of the Directory
func Path(d Directory) string {
	var path string
	for _, node := range d.Ancestry() {
		path += "/" + node.Name()
	}
	return path
}

// FilePath returns the full path of the File
func FilePath(f File) string {
	if f.Directory() != nil {
		return Path(f.Directory()) + "/" + f.Name()
	}
	return f.Name()
}

// VersionPath returns the full path of the Version
func VersionPath(v Version) string {
	if v.File() != nil {
		return FilePath(v.File()) + "/" + v.UUID().String()
	}
	return v.UUID().String()
}
