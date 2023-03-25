package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id   uuid.UUID
	name string
}

func MakeUser(id uuid.UUID, name string) User {
	return User{id: id, name: name}
}

type UserEmailPassword struct {
	user     User
	email    string
	password Password
}

func (u UserEmailPassword) User() User {
	return u.user
}

func (u UserEmailPassword) Email() string {
	return u.email
}

func (u UserEmailPassword) HashPassword() string {
	return u.password.Hash()
}

func NewUserWithEmailPassword(name string, email string, password Password) (UserEmailPassword, error) {
	user := User{
		id:   uuid.New(),
		name: name,
	}
	return UserEmailPassword{user: user, email: email, password: password}, nil
}

func (u User) ID() uuid.UUID {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type Email string

func MakeEmail(value string) Email {
	return Email(value)
}

func (e Email) String() string {
	return string(e)
}

type Password struct {
	value  string
	hashed string
}

func MakePassword(value string) (Password, error) {
	hashed, err := hashPassword(value)
	if err != nil {
		return Password{}, err
	}

	return Password{value: value, hashed: hashed}, nil
}

func (p Password) String() string {
	return p.value
}

func (p Password) Hash() string {
	return p.hashed
}
