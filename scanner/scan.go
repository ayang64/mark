package scanner

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ayang64/mark/token"
)

type Scanner struct {
	rs io.RuneScanner
}

// should only be called after a successful .Peek()
func (s *Scanner) Consume() rune {
	r, _, err := s.rs.ReadRune()
	if err != nil {
		return utf8.RuneError
	}
	return r
}

func (s *Scanner) Peek() (rune, error) {
	r, _, err := s.rs.ReadRune()
	if err == nil {
		s.rs.UnreadRune()
	}
	return r, err
}

// tokens can a string representation
type Token interface {
	String() string
}

func match(rs io.RuneScanner, match func(r rune) (bool, bool, error)) string {
	sb := strings.Builder{}
mainloop:
	for {
		r, _, err := rs.ReadRune()
		if err != nil {
			break
		}
		accept, cont, err := match(r)
		switch {
		case err != nil:
			rs.UnreadRune()
			break mainloop
		case accept:
			sb.WriteRune(r)
		}
		if !cont {
			break mainloop
		}
	}
	return sb.String()
}

//
//  = title   h1
//  == sub section h2
//  === hN where N is the number of equals
//
//  *bold text*
//  /italicized text here/
//  _underline this text_
//
//  [space]* bulleted item
//
//  [url](text)  hyper link
//
//  ```
//  code block
//  ```
//
//  `preformatted text`
//
//  ![url.jpeg](desc) image
//

// scans next token
func (s *Scanner) Next() (token.Token, error) {
	r, err := s.Peek()
	if err != nil {
		return token.Error{Err: err}, err
	}

	switch {
	case r == '\\':
		// an escaped character.
		s.Consume() // consume current rune
		v := match(s.rs, func(r rune) (bool, bool, error) {
			return true, false, nil
		})
		return token.Rune([]rune(v)[0]), nil
	case r == '\n':
		s.Consume()
		return token.Rune(r), nil
	case r == '*', r == '_', r == '/':
		return token.Rune(s.Consume()), nil
	case unicode.IsSpace(r):
		v := match(s.rs, func(r rune) (bool, bool, error) {
			if unicode.IsSpace(r) {
				return true, true, nil
			}
			return false, false, fmt.Errorf("%c is not a space rune", r)
		})
		return token.Space(v), nil
	default:
		// read until a space rune is detected.
		v := match(s.rs, func(r rune) (bool, bool, error) {
			if unicode.IsSpace(r) {
				return false, false, fmt.Errorf("%v is not in class ATOM", r)
			}
			return true, true, nil
		})
		return token.Atom(v), nil
	}
}
