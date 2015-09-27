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
	DetachDirectory(Directory) error
	AttachFile(File) error
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
	SetDirectory(Directory) error
	CurrentVersion() Version
	SetCurrentVersion(Version) error
	Owner() User
	SetOwner(User) error
	Delete() error
	AttachVersion(Version) error
	DetachVersion(Version) error
	FindVersions(time.Time, time.Time, User, int) ([]Version, error)
}

// Version is a version
type Version interface {
	UUID() uuid.UUID
	File() File
	SetFile(File) error
	Time() time.Time
	SetTime(time.Time) error
	Creator() User
	SetCreator(User) error
	Delete() error
}

// User is a user
type User interface {
	UUID() uuid.UUID
}
