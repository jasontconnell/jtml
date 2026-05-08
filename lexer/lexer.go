package lexer

import (
	"bufio"
	"bytes"
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
	buf := bytes.NewBufferString(l.input)
	scn := bufio.NewScanner(buf)
	tokens := []Token{}
	linenum := 0
	level := 0
	curstate := None

	for scn.Scan() {
		line := scn.Text()
		incline := false
		prefix := true

		for i := 0; i < len(line); i++ {
			switch line[i] {
			case ' ':
				if prefix {
					level++
				}
			case '@', '#':
				sp := i
				if prefix {
					sp = i + 1
					curstate = Directive
					if line[i] == '#' {
						curstate = Include
					}
					incline = true
				} else {
					curstate = Raw
					incline = false
				}

				identifier, endline := l.getIdentifier(line, sp)
				tk := Token{
					Type:    curstate,
					Start:   i,
					Level:   level,
					LineNum: linenum,
					Value:   identifier,
					Endline: endline,
				}
				tokens = append(tokens, tk)
				i += len(identifier)

				prefix = false
			default:
				identifier, endline := l.getIdentifier(line, i)
				idtype := Raw
				if incline {
					idtype = Parameter
				}
				tk := Token{
					Type:    idtype,
					Start:   i,
					Value:   identifier,
					LineNum: linenum,
					Level:   level,
					Endline: endline,
				}
				tokens = append(tokens, tk)
				i += len(identifier)
				prefix = false
			}
		}

		linenum++
		level = 0
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
