package dir

import (
	"regexp"
	"sync"

	"github.com/dustywilson/dir"
	"github.com/satori/go.uuid"
)

// Directory is an implementation of dir.Directory
type Directory struct {
	root     dir.Directory
	uuid     uuid.UUID
	name     string
	parent   dir.Directory
	children []dir.Directory
	files    []dir.File
	owner    dir.User
	sync.RWMutex
}

// NewDirectory returns a new Directory
func NewDirectory() dir.Directory {
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
func (d *Directory) Parent() dir.Directory {
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
func (d *Directory) Root() dir.Directory {
	d.RLock()
	defer d.RUnlock()
	return d.root
}

// SetRoot sets/replaces the root dir.Directory for this Directory
// This must only be called from method PlaceDirectory() on an implementation of dir.Directory.
func (d *Directory) SetRoot(root dir.Directory) error {
	d.Lock()
	defer d.Unlock()
	d.root = root
	return nil
}

// Owner returns the owner's User
func (d *Directory) Owner() dir.User {
	d.RLock()
	defer d.RUnlock()
	return d.owner
}

// SetOwner changes the Directory's owner
func (d *Directory) SetOwner(owner dir.User) error {
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
		return dir.ErrNotEmpty
	}
	if d.parent != nil {
		d.parent.DetachDirectory(d)
	} else {
		return dir.ErrIsRoot
	}
	return nil
}

// AttachDirectory adds a dir.Directory to this Directory
// You likely would use the .Create() or .CreateDirectory() method instead, if available.
func (d *Directory) AttachDirectory(dirIn dir.Directory) error {
	d.Lock()
	defer d.Unlock()
	dirIn.SetRoot(d.root)
	// TODO: test for errors; ensure no duplicate attachment
	d.children = append(d.children, dirIn)
	return nil
}

// DetachDirectory removes a dir.Directory from this Directory
// This is intended to be called from the target Direcotry's .Delete() method itself.
// This does not test that the Directory is ready to be deleted, whatever that means (is empty, etc).
func (d *Directory) DetachDirectory(dirIn dir.Directory) error {
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
		return dir.ErrNoMatch
	}
	return nil
}

// AttachFile adds a dir.File to this Directory
// If your File has a .Create() method, you might use that instead.
func (d *Directory) AttachFile(f dir.File) error {
	d.Lock()
	defer d.Unlock()
	// TODO: test for errors; ensure no duplicate attachment
	d.files = append(d.files, f)
	return nil
}

// DetachFile removes a dir.File from this Directory
// This is intended to be called from the target File's .Delete() method itself.
// This does not test that the File is ready to be deleted, whatever that means.
func (d *Directory) DetachFile(f dir.File) error {
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
		return dir.ErrNoMatch
	}
	return nil
}

// FindDirectories searches for one or more dir.Directory entries, recursively starting at this Directory
func (d *Directory) FindDirectories(search *regexp.Regexp, recurseLevel int) ([]dir.Directory, error) {
	d.RLock()
	defer d.RUnlock()
	var directories []dir.Directory
	if search.MatchString(d.name) {
		directories = append(directories, d)
	}
	for _, dir := range d.children {
		if recurseLevel != 0 {
			childMatches, err := dir.FindDirectories(search, recurseLevel-1)
			if err != nil {
				// TODO: decide what to do here, if anything.
			} else {
				directories = append(directories, childMatches...)
			}
		}
	}
	if len(directories) == 0 {
		return nil, dir.ErrNoMatch
	}
	return directories, nil
}

// FindFiles searches for one or more File entries, recursively starting at the Directory
func (d *Directory) FindFiles(search *regexp.Regexp, recurseLevel int) ([]dir.File, error) {
	d.RLock()
	defer d.RUnlock()
	var files []dir.File
	for _, file := range d.files {
		if search.MatchString(file.Name()) {
			files = append(files, file)
		}
	}
	for _, dir := range d.children {
		if recurseLevel != 0 {
			childMatches, err := dir.FindFiles(search, recurseLevel-1)
			if err != nil {
				// TODO: decide what to do here, if anything.
			} else {
				files = append(files, childMatches...)
			}
		}
	}
	if len(files) == 0 {
		return nil, dir.ErrNoMatch
	}
	return nil, nil
}

// CreateDirectory creates a sub-Directory of the provided Directory
func CreateDirectory(d dir.Directory, name string) (dir.Directory, error) {
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
