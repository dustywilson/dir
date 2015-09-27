package dir

import (
	"regexp"
	"time"

	"github.com/satori/go.uuid"
)

// Directory is a directory
type Directory interface {
	UUID() uuid.UUID
	Name() string
	Rename(string) error
	Parent() Directory
	Ancestry() []Directory
	IsRoot() bool
	Root() Directory
	SetRoot(Directory) error
	Owner() User
	SetOwner(User) error
	Delete() error
	AttachDirectory(Directory) error
	AttachFile(File) error
	DetachDirectory(Directory) error
	DetachFile(File) error
	FindDirectories(*regexp.Regexp, int) ([]Directory, error)
	FindFiles(*regexp.Regexp, int) ([]File, error)
}

// File is a file
type File interface {
	UUID() uuid.UUID
	Name() string
	Rename(string) error
	Directory() Directory
	CurrentVersion() Version
	SetCurrentVersion(Version) error
	Owner() User
	SetOwner(User) error
	Delete() error
	AttachVersion(Version) error
	FindVersions(time.Time, time.Time, User, int) ([]Version, error)
}

// Version is a version
type Version interface {
	UUID() uuid.UUID
	File() File
	Time() time.Time
	Creator() User
	Delete() error
	PlaceVersion(File) error
}

// User is a user
type User interface {
	UUID() uuid.UUID
}
