package dir

// Path returns the full path of the Directory
func Path(d Directory) string {
	var path string
	for _, node := range d.Ancestry() {
		path += "/" + node.Name()
	}
	return path
}
