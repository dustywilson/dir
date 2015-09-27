package user

import (
	"sync"

	"github.com/dustywilson/dir"
	"github.com/satori/go.uuid"
)

// User is an implementation of dir2.User
type User struct {
	uuid uuid.UUID
	name string
	sync.RWMutex
}

// NewUser returns a new User
func NewUser(name string) dir.User {
	return &User{
		uuid: uuid.NewV4(),
		name: name,
	}
}

func (u *User) String() string {
	return u.Name()
}

// UUID returns the User's UUID
func (u *User) UUID() uuid.UUID {
	u.RLock()
	defer u.RUnlock()
	return u.uuid
}

// Name returns the User's name
func (u *User) Name() string {
	u.RLock()
	defer u.RUnlock()
	return u.name
}

// SetName sets the User's name
func (u *User) SetName(name string) error {
	u.Lock()
	defer u.Unlock()
	if len(name) == 0 {
		return dir.ErrIsEmpty
	}
	u.name = name
	// TODO: check for collisions or validity or something?
	return nil
}
