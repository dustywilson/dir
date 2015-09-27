package dir

import (
	"errors"
	"regexp"
	"time"

	"github.com/satori/go.uuid"
)

// Directory is a directory
type Directory interface {
	UUID() uuid.UUID
	Name() string
	Ancestry() []Directory
	Rename(string) error
	Parent() Directory
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
	Owner() User
	Delete() error
	PlaceFile(Directory) error
	PlaceVersion(Version) error
	FindVersions(string) ([]Version, error)
}

// Version is a version
type Version interface {
	UUID() uuid.UUID
	File() File
	Timestamp() time.Time
	Creator() User
	Delete() error
	PlaceVersion(File) error
}

// User is a user
type User interface {
	UUID() uuid.UUID
}

// Errors
var (
	ErrIsRoot   = errors.New("is root")
	ErrNoMatch  = errors.New("no match")
	ErrNotEmpty = errors.New("not empty")
)
