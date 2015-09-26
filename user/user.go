package user

import "github.com/satori/go.uuid"

// User is an implementation of dir2.User
type User struct {
	uuid uuid.UUID
}

// NewUser returns a new User
func NewUser() User {
	return User{
		uuid: uuid.NewV4(),
	}
}

// UUID returns the User's UUID
func (u User) UUID() uuid.UUID {
	return u.uuid
}
