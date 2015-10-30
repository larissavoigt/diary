package user

import "strings"

type User struct {
	ID         int64
	FacebookID string
	Token      string
	Name       string
}

func (u *User) FirstName() string {
	s := strings.Split(u.Name, " ")
	return s[0]
}
