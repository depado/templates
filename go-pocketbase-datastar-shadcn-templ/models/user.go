package models

import (
	"strings"
	"unicode/utf8"

	"github.com/pocketbase/pocketbase/core"
)

type User struct {
	*core.Record
}

func (u *User) Email() string { return u.GetString("email") }

func (u *User) Name() string { return u.GetString("name") }

func (u *User) Avatar() string { return u.GetString("avatar") }

func (u *User) Initials() string {
	name := strings.TrimSpace(u.Name())
	if name == "" {
		return "?"
	}
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	first, _ := utf8.DecodeRuneInString(parts[0])
	initials := string(first)
	if len(parts) > 1 {
		last, _ := utf8.DecodeRuneInString(parts[len(parts)-1])
		initials += string(last)
	}
	return strings.ToUpper(initials)
}

func NewUser(record *core.Record) *User {
	return &User{Record: record}
}
