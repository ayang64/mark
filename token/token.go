package token

import "fmt"

type Token interface {
	String() string
}

type Rune rune

func (r Rune) String() string {
	return fmt.Sprintf("%c", r)
}

type Error struct{ Err error }

func (e Error) String() string {
	return e.Err.Error()
}

type EOL struct{}

func (EOL) String() string {
	return "<EOL>"
}

type BOL struct{}

func (BOL) String() string {
	return "<BOL>"
}

type Space string

func (s Space) String() string {
	return fmt.Sprintf("%q", string(s))
}

type Atom string

func (a Atom) String() string {
	return string(a)
}
