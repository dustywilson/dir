package version

import (
	"sync"
	"time"

	"github.com/dustywilson/dir"
	"github.com/satori/go.uuid"
)

// Version is an implementation of dir.Version
type Version struct {
	uuid    uuid.UUID
	file    dir.File
	creator dir.User
	time    time.Time
	sync.RWMutex
}

func (v *Version) String() string {
	return dir.VersionPath(v)
}

// UUID returns the version's UUID
func (v *Version) UUID() uuid.UUID {
	v.RLock()
	defer v.RUnlock()
	return v.uuid
}

// File returns the version's File
func (v *Version) File() dir.File {
	v.RLock()
	defer v.RUnlock()
	return v.file
}

// SetFile sets the version's File
func (v *Version) SetFile(file dir.File) error {
	v.Lock()
	defer v.Unlock()
	v.file = file
	// TODO: check for some sort of error?
	return nil
}

// Time returns the version's time.Time
func (v *Version) Time() time.Time {
	v.RLock()
	defer v.RUnlock()
	return v.time
}

// SetTime sets the version's time.Time
func (v *Version) SetTime(t time.Time) error {
	v.Lock()
	defer v.Unlock()
	v.time = t
	// TODO: check for some sort of error?
	return nil
}

// Creator returns the version's creator dir.User
func (v *Version) Creator() dir.User {
	v.RLock()
	defer v.RUnlock()
	return v.creator
}

// SetCreator sets the version's dir.User
func (v *Version) SetCreator(u dir.User) error {
	v.Lock()
	defer v.Unlock()
	v.creator = u
	// TODO: check for some sort of error?
	return nil
}

// Delete deletes the Version
func (v *Version) Delete() error {
	v.Lock()
	defer v.Unlock()
	if v.file != nil {
		v.file.DetachVersion(v)
	}
	// TODO: should we set some sort of "isDeleted" flag?
	return nil
}
