package directory

import (
	"regexp"
	"sync"

	"github.com/dustywilson/directory"
	"github.com/satori/go.uuid"
)

// Directory is an implementation of directory.Directory
type Directory struct {
	root     directory.Directory
	uuid     uuid.UUID
	name     string
	parent   directory.Directory
	children []directory.Directory
	files    []directory.File
	owner    directory.User
	sync.RWMutex
}

// NewDirectory returns a new Directory
func NewDirectory() directory.Directory {
	d := &Directory{
		uuid: uuid.NewV4(),
	}
	d.root = d
	return d
}

// UUID returns the Directory's UUID
func (d *Directory) UUID() uuid.UUID {
	return d.uuid
}

// Name returns the Directory's name
func (d *Directory) Name() string {
	d.RLock()
	defer d.RUnlock()
	return d.name
}

// Rename changes the Directory's name
func (d *Directory) Rename(name string) error {
	d.Lock()
	defer d.Unlock()
	d.name = name
	// TODO: check for name collision, nil Directory, etc
	return nil
}

// Parent returns the parent's Directory
func (d *Directory) Parent() directory.Directory {
	d.RLock()
	defer d.RUnlock()
	return d.parent
}

// IsRoot returns bool if this Directory is the root
func (d *Directory) IsRoot() bool {
	d.RLock()
	defer d.RUnlock()
	return d.parent == nil
}

// Root returns the root Directory
func (d *Directory) Root() directory.Directory {
	d.RLock()
	defer d.RUnlock()
	return d.root
}

// SetRoot sets/replaces the root directory.Directory for this Directory
// This must only be called from method PlaceDirectory() on an implementation of directory.Directory.
func (d *Directory) SetRoot(root directory.Directory) error {
	d.Lock()
	defer d.Unlock()
	d.root = root
	return nil
}

// Owner returns the owner's User
func (d *Directory) Owner() directory.User {
	d.RLock()
	defer d.RUnlock()
	return d.owner
}

// SetOwner changes the Directory's owner
func (d *Directory) SetOwner(owner directory.User) error {
	d.Lock()
	defer d.Unlock()
	d.owner = owner
	// TODO: check for some sort of error?
	return nil
}

// Delete deletes the Directory
func (d *Directory) Delete() error {
	d.Lock()
	defer d.Unlock()
	if len(d.children) > 0 || len(d.files) > 0 {
		return directory.ErrNotEmpty
	}
	if d.parent != nil {
		d.parent.DetachDirectory(d)
	} else {
		return directory.ErrIsRoot
	}
	return nil
}

// AttachDirectory adds a directory.Directory to this Directory
// You likely would use the .Create() or .CreateDirectory() method instead, if available.
func (d *Directory) AttachDirectory(dirIn directory.Directory) error {
	d.Lock()
	defer d.Unlock()
	dirIn.SetRoot(d.root)
	// TODO: test for errors; ensure no duplicate attachment
	d.children = append(d.children, dirIn)
	return nil
}

// DetachDirectory removes a directory.Directory from this Directory
// This is intended to be called from the target Direcotry's .Delete() method itself.
// This does not test that the Directory is ready to be deleted, whatever that means (is empty, etc).
func (d *Directory) DetachDirectory(dirIn directory.Directory) error {
	d.Lock()
	defer d.Unlock()
	matched := false
	for i, dir := range d.children {
		if dir == dirIn {
			d.children = append(d.children[:i], d.children[i+1:]...)
			matched = true
		}
	}
	if !matched {
		return directory.ErrNoMatch
	}
	return nil
}

// AttachFile adds a directory.File to this Directory
// If your File has a .Create() method, you might use that instead.
func (d *Directory) AttachFile(f directory.File) error {
	d.Lock()
	defer d.Unlock()
	// TODO: test for errors; ensure no duplicate attachment
	d.files = append(d.files, f)
	return nil
}

// DetachFile removes a directory.File from this Directory
// This is intended to be called from the target File's .Delete() method itself.
// This does not test that the File is ready to be deleted, whatever that means.
func (d *Directory) DetachFile(f directory.File) error {
	d.Lock()
	defer d.Unlock()
	matched := false
	for i, file := range d.files {
		if file == f {
			d.files = append(d.files[:i], d.files[i+1:]...)
			matched = true
		}
	}
	if !matched {
		return directory.ErrNoMatch
	}
	return nil
}

// FindDirectories searches for one or more directory.Directory entries, recursively starting at this Directory
func (d *Directory) FindDirectories(search *regexp.Regexp, recurseLevel int) ([]directory.Directory, error) {
	d.RLock()
	defer d.RUnlock()
	var directories []directory.Directory
	if search.MatchString(d.name) {
		directories = append(directories, d)
	}
	for _, directory := range d.children {
		if recurseLevel != 0 {
			childMatches, err := directory.FindDirectories(search, recurseLevel-1)
			if err != nil {
				// TODO: decide what to do here, if anything.
			} else {
				directories = append(directories, childMatches...)
			}
		}
	}
	if len(directories) == 0 {
		return nil, directory.ErrNoMatch
	}
	return directories, nil
}

// FindFiles searches for one or more File entries, recursively starting at the Directory
func (d *Directory) FindFiles(search *regexp.Regexp, recurseLevel int) ([]directory.File, error) {
	d.RLock()
	defer d.RUnlock()
	var files []directory.File
	for _, file := range d.files {
		if search.MatchString(file.Name()) {
			files = append(files, file)
		}
	}
	for _, directory := range d.children {
		if recurseLevel != 0 {
			childMatches, err := directory.FindFiles(search, recurseLevel-1)
			if err != nil {
				// TODO: decide what to do here, if anything.
			} else {
				files = append(files, childMatches...)
			}
		}
	}
	if len(files) == 0 {
		return nil, directory.ErrNoMatch
	}
	return nil, nil
}

// CreateDirectory creates a sub-Directory of the provided Directory
func CreateDirectory(d directory.Directory, name string) (directory.Directory, error) {
	subdir := &Directory{
		uuid: uuid.NewV4(),
		name: name,
	}
	if d != nil {
		subdir.root = d.Root()
		subdir.parent = d
		err := d.AttachDirectory(subdir)
		if err != nil {
			return nil, err
		}
	} else {
		subdir.root = subdir
	}
	return subdir, nil
}
