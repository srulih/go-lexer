package main

import (
	"fmt"
	"regexp"
	"strings"
)


func MapValues(m map[string]string) []string {
	v := make([]string, len(m), len(m))
    idx := 0
    for  key, _ := range m {
       v[idx] = key
       idx++
    }
    return v
}

type Token struct {
	Pos int
	Value interface{}
	Type string
}

func (t Token) String() string {
	return fmt.Sprintf("%s with value %s at position %d", t.Type, t.Value, t.Pos)
}

type Lexer struct {
	rules map[string]string
	regex *regexp.Regexp
	whitespaceRegex *regexp.Regexp
	skipWhitespace bool
	buffer string
	pos int
}

func BuildLexer(rules map[string]string, skipWhitespace bool) Lexer {
	rs := strings.Join(MapValues(rules), "|")
	r := regexp.MustCompile(rs)
	rw, _ := regexp.Compile(`\S`)
	return Lexer{
		skipWhitespace: skipWhitespace,
		rules: rules,
		regex: r,
		whitespaceRegex: rw,
	}
}

func (l *Lexer) Input(buf string) {
	l.buffer = buf
}

func (l *Lexer) Token() (*Token, error) {
	if l.pos >= len(l.buffer) {
		return nil, nil
	}
	if l.skipWhitespace {
		m := l.whitespaceRegex.FindAllStringIndex(l.buffer[l.pos:],1)
		if len(m) > 0 {
			l.pos += m[0][0]
		}
	}
	r := l.regex.FindAllStringIndex(l.buffer[l.pos:], 1)
	if len(r) > 0 && r[0][0] == 0 {
		value := l.buffer[l.pos: l.pos + r[0][1]]
		var ty string
		for k,v := range l.rules {
			if m, _:= regexp.MatchString(k, value); m {
				ty = v
			}
		}
		t := Token{
			Value: value,
			Pos: l.pos,
			Type: ty, 
		}
		l.pos += r[0][1]
		return &t, nil
	}
	return nil, fmt.Errorf("could not match anything at position %d", l.pos)
}

func (l *Lexer) Tokens() []*Token {
	tokens := make([]*Token, 0)
	for {
		tok, err := l.Token()
		if err != nil {
			panic(err)
		}
		if tok == nil {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

func main() {
	m := map[string]string {
		`\d+`: "NUMBER",
		`[a-zA-Z_]\w+`: "IDENTIFIER",
		`\+`: "PLUS",
	}
	l := BuildLexer(m, true)
	l.Input("hello 23 + 34    ft")
	toks := l.Tokens()
	for _, v := range toks {
		fmt.Println(v)
	}
}

