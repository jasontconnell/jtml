package lexer

import (
	"fmt"
	"strings"
)

type lexer struct {
	input string
	pos   int
}

type State int

const (
	None State = iota
	Include
	Directive
	Parameter
	Raw
)

func (s State) String() string {
	str := "None"
	switch s {
	case Include:
		str = "Include"
	case Directive:
		str = "Directive"
	case Parameter:
		str = "Parameter"
	case Raw:
		str = "Raw"
	}
	return str
}

type Token struct {
	Value   string
	Start   int
	Type    State
	Level   int
	LineNum int
	Endline bool
}

func (t Token) String() string {
	return fmt.Sprintf("TOKEN[%s start: %d type: %v level: %d line: %d endline: %v]", t.Value, t.Start, t.Type, t.Level, t.LineNum, t.Endline)
}

func Lex(input string) []Token {
	lex := &lexer{input: input}
	return lex.getTokens()
}

func (l *lexer) getTokens() []Token {
	tokens := []Token{}
	linenum := 0
	level := 0
	prefix := true
	curstate := None
	incline := false
	for pos := 0; pos < len(l.input); pos++ {
		switch l.input[pos] {
		case ' ':
			if prefix {
				level++
			}
		case '@', '#':
			sp := pos
			if prefix {
				sp = pos + 1
				curstate = Directive
				if l.input[pos] == '#' {
					curstate = Include
				}
				incline = true
			} else {
				curstate = Raw
				incline = false
			}

			identifier, endline := l.getIdentifier(l.input, sp)
			tk := Token{
				Type:    curstate,
				Start:   pos,
				Level:   level,
				LineNum: linenum,
				Value:   identifier,
				Endline: endline,
			}
			tokens = append(tokens, tk)
			pos += len(identifier)

			prefix = false
		case '\n':
			linenum++
			prefix = true
			incline = false
			level = 0
		case '\r':
			prefix = true
			incline = false
			continue
		default:
			identifier, endline := l.getIdentifier(l.input, pos)
			idtype := Raw
			if incline {
				idtype = Parameter
			}
			tk := Token{
				Type:    idtype,
				Start:   pos,
				Value:   identifier,
				LineNum: linenum,
				Level:   level,
				Endline: endline,
			}
			tokens = append(tokens, tk)
			pos += len(identifier)
			prefix = false
		}
	}
	return tokens
}

func (l *lexer) getIdentifier(input string, pos int) (string, bool) {
	id := ""
	endline := false
	for i := pos; i < len(input); i++ {
		if input[i] == ' ' {
			break
		}
		if input[i] == '\n' || input[i] == '\r' {
			endline = true
			break
		}
		id += strings.TrimSpace(string(input[i]))
		if i == len(input)-1 {
			endline = true
		}
	}
	return id, endline
}
