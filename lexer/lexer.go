package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"slices"
	"sort"
	"strings"
)

type State int

const (
	None State = iota
	Include
	Directive
	Parameter
	Raw
	Endline
	Indent
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
	case Endline:
		str = "Endline"
	case Indent:
		str = "Indent"
	}
	return str
}

type Token struct {
	Value   string
	Start   int
	Type    State
	LineNum int
}

func (t Token) String() string {
	return fmt.Sprintf("TOKEN[%s start: %d type: %v line: %d]", t.Value, t.Start, t.Type, t.LineNum)
}

func Lex(input string) []Token {
	lines := normalizeInput(input)
	return getTokens(lines)
}

type spaceres struct {
	line   string
	spaces int
}

func normalizeInput(s string) []string {
	buf := bytes.NewBufferString(s)
	scn := bufio.NewScanner(buf)

	lens := []int{}
	spacereses := []spaceres{}
	v := make(map[int]int)

	for scn.Scan() {
		line := scn.Text()

		if len(strings.Trim(line, " ")) == 0 {
			continue
		}

		tmp := strings.TrimLeft(line, " ")
		spaces := len(line) - len(tmp)

		if _, ok := v[spaces]; !ok {
			lens = append(lens, spaces)
		}
		v[spaces] = spaces
		spacereses = append(spacereses, spaceres{strings.TrimRight(tmp, "\r\n"), spaces})
	}

	sort.Ints(lens)

	lines := []string{}
	for _, res := range spacereses {
		level := slices.Index(lens, res.spaces)
		if level == -1 {
			continue
		}

		lines = append(lines, strings.Repeat(" ", level)+res.line)
	}

	return lines
}

func getTokens(lines []string) []Token {
	tokens := []Token{}
	level := 0
	curstate := None

	for linenum, line := range lines {
		incline := false
		prefix := true
		comment := false

		for i := 0; i < len(line) && !comment; i++ {
			switch line[i] {
			case ' ':
				if prefix {
					tokens = append(tokens, Token{Value: " ", Start: i, Type: Indent, LineNum: linenum})
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

				identifier := getIdentifier(line, sp)
				tk := Token{
					Type:    curstate,
					Start:   i,
					LineNum: linenum,
					Value:   identifier,
				}
				tokens = append(tokens, tk)
				i += len(identifier)

				prefix = false
			case '`':
				if prefix {
					comment = true
				}
			default:
				var tk Token
				identifier := getIdentifier(line, i)
				idtype := Raw
				if incline && line[i] == '[' {
					idtype = Parameter
					val := getParamString(line, i+1)
					tk = Token{
						Type:    Parameter,
						Value:   val,
						LineNum: linenum,
						Start:   i,
					}
					i += len(val) + 1
				} else {
					tk = Token{
						Type:    idtype,
						Start:   i,
						Value:   identifier,
						LineNum: linenum,
					}
					i += len(identifier)
				}

				tokens = append(tokens, tk)
				prefix = false
			}
		}

		el := Token{
			Type:    Endline,
			Start:   len(line) - 1,
			LineNum: linenum,
			Value:   "_endline_",
		}
		tokens = append(tokens, el)

		level = 0
	}

	return tokens
}

func getIdentifier(input string, pos int) string {
	id := ""
	for i := pos; i < len(input); i++ {
		if input[i] == ' ' {
			break
		}
		id += strings.TrimSpace(string(input[i]))
	}
	return id
}

func getParamString(input string, pos int) string {
	str := ""
	for i := pos; i < len(input); i++ {
		if input[i] == ']' {
			break
		}
		str += string(input[i])
	}
	return str
}
