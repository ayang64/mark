package scanner

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"unicode"
)

func TestScanner(t *testing.T) {
	data := `= Hello, World
` + "```" + `
block quote
` + "```" + `
* a1 \* \\ \X
  * b1 
* a2
	* b1
	* b2
		* c1

== second level heading

paragraph text
`
	s := Scanner{
		rs: strings.NewReader(data),
	}

	for {
		t, err := s.Next()
		if err != nil {
			break
		}
		log.Printf("%[1]s (%[1]T)", t)
	}

	s.rs = strings.NewReader(data)

	for t := range s.Scan() {
		log.Printf("sc: %[1]s (%[1]T)", t)
	}
}

func TestMatch(t *testing.T) {
	scan := func(s string) []string {
		rs := strings.NewReader(s)
		m := []string(nil)
		for {
			v, _, err := rs.ReadRune()
			if err != nil {
				break
			}
			rs.UnreadRune()

			switch {
			case unicode.IsSpace(v):
				match(rs, func(r rune) (bool, bool, error) {
					if !unicode.IsSpace(r) {
						return false, false, fmt.Errorf("not a space")
					}
					return true, true, nil
				})
			default:
				v := match(rs, func(r rune) (bool, bool, error) {
					if unicode.IsSpace(r) {
						return false, false, fmt.Errorf("not a space")
					}
					return true, true, nil
				})
				m = append(m, v)

			}
		}
		return m
	}

	t.Logf("%#v", scan("123 456 789 101112 a b c d e f g"))
}
