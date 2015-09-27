package file

import (
	"sync"
	"time"

	"github.com/dustywilson/dir"
	"github.com/satori/go.uuid"
)

// File is an implementation of dir.File
type File struct {
	uuid           uuid.UUID
	name           string
	directory      dir.Directory
	versions       []dir.Version
	currentVersion dir.Version
	owner          dir.User
	sync.RWMutex
}

// UUID returns the file's UUID
func (f *File) UUID() uuid.UUID {
	f.RLock()
	defer f.RUnlock()
	return f.uuid
}

// Name returns the file's name
func (f *File) Name() string {
	f.RLock()
	defer f.RUnlock()
	// TODO: check for name collision, nil File, etc
	return f.name
}

// Rename sets the file's name
func (f *File) Rename(name string) error {
	f.Lock()
	defer f.Unlock()
	f.name = name
	return nil
}

// Directory returns the File's Directory (if attached to one)
func (f *File) Directory() dir.Directory {
	f.RLock()
	defer f.RUnlock()
	return f.directory
}

// SetDirectory sets the File's Directory
func (f *File) SetDirectory(d dir.Directory) error {
	f.Lock()
	defer f.Unlock()
	f.directory = d
	// TODO: check for some sort of error?
	return nil
}

// CurrentVersion returns the File's current Version (which might not be the newest one)
func (f *File) CurrentVersion() dir.Version {
	f.RLock()
	defer f.RUnlock()
	return f.currentVersion
}

// SetCurrentVersion sets the File's current Version (doesn't have to be the newest Version)
func (f *File) SetCurrentVersion(currentVersion dir.Version) error {
	f.Lock()
	defer f.Unlock()
	f.attachVersionLocked(currentVersion) // noop if already attached
	f.currentVersion = currentVersion
	return nil
}

// Owner returns the file owner's dir.User
func (f *File) Owner() dir.User {
	f.RLock()
	defer f.RUnlock()
	return f.owner
}

// SetOwner returns the file owner's dir.User
func (f *File) SetOwner(owner dir.User) error {
	f.Lock()
	defer f.Unlock()
	f.owner = owner
	// TODO: check for some sort of error?
	return nil
}

// Delete deletes the File
func (f *File) Delete() error {
	f.Lock()
	defer f.Unlock()
	if len(f.versions) > 0 {
		// we can't delete files that have attached versions
		return dir.ErrNotEmpty
	}
	if f.directory != nil {
		f.directory.DetachFile(f)
	}
	// TODO: should we set some sort of "isDeleted" flag?
	return nil
}

// AttachVersion attaches a dir.Version to this file
// This checks to ensure the provided dir.Version isn't a duplicate before it attaches it.
func (f *File) AttachVersion(v dir.Version) error {
	f.Lock()
	defer f.Unlock()
	return f.attachVersionLocked(v)
}

func (f *File) attachVersionLocked(v dir.Version) error {
	// Warning!  This assumes the file is already locked for write by the caller.
	for _, ver := range f.versions {
		if ver == v {
			return dir.ErrExists
		}
	}
	f.versions = append(f.versions, v)
	return nil
}

// DetachVersion removes a dir.Version from this File
// This is intended to be called from the target File's .Delete() method itself.
// This does not test that the Version is ready to be deleted, whatever that means.
func (f *File) DetachVersion(v dir.Version) error {
	f.Lock()
	defer f.Unlock()
	if f.currentVersion == v {
		if len(f.versions) > 1 {
			// can't delete the current version.  set the currentVersion to something else first.
			return dir.ErrIsCurrentVersion
		}
		f.currentVersion = nil // because we're about to delete the _only_ version, we'll nil this.
	}
	matched := false
	for i, ver := range f.versions {
		if ver == v {
			f.versions = append(f.versions[:i], f.versions[i+1:]...)
			matched = true
		}
	}
	if !matched {
		return dir.ErrNoMatch
	}
	return nil
}

// FindVersions searches for one or more Version entries
// The `after`, `before`, and `creator` arguments are all optional and will be used in an `AND` search, if provided.
func (f *File) FindVersions(after time.Time, before time.Time, creator dir.User, recurseLevel int) ([]dir.Version, error) {
	f.RLock()
	defer f.RUnlock()
	var matches []dir.Version
	for _, v := range f.versions {
		if (after.IsZero() || after.After(v.Time())) && (before.IsZero() || !before.After(v.Time())) && (creator == nil || creator == v.Creator()) {
			matches = append(matches, v)
		}
	}
	return matches, nil
}
