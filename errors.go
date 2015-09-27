package dir

import "errors"

// Errors
var (
	ErrIsRoot           = errors.New("is root")
	ErrNoMatch          = errors.New("no match")
	ErrNotEmpty         = errors.New("not empty")
	ErrExists           = errors.New("already exists")
	ErrIsCurrentVersion = errors.New("is set as current version")
)
