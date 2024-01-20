package main

import (
    "io"
    "os"
    "fmt"
    "strings"
)

func readFile(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
	return "", err
    }

    fileSize, err := file.Seek(0, 2)
    if err != nil {
	return "", err
    }

    _, err = file.Seek(0, 0)
    if err != nil {
	return "", err
    }

    buffer := make([]byte, fileSize)

    _, err = file.Read(buffer)
    if err != nil && err != io.EOF {
	return "", err
    }

    return string(buffer), nil
}

type TokenKind int

const (
    TokenUnknown TokenKind = iota
    TokenHeading
    TokenParagraph
    TokenCount
)

func (kind TokenKind) toString() string {
    if TokenCount > 3 {
	panic("missing name of new token added")
    }
    switch kind {
	case TokenHeading: return "TokenHeading"
	case TokenParagraph: return "TokenParagraph"
    }
    return "TokenUnknown" 
}

type Token struct {
    kind  TokenKind
    level int
    value string
}

func (t *Token) toString() string {
    return fmt.Sprintf("Token {%s, %d, '%s'}\n", t.kind.toString(), t.level, t.value)
}

func parseHeading(line string) Token {
    i := 0
    for ; i < len(line); i += 1 {
	if line[i] != '#' {
	    break
	}
    }
    if i > 6 {
	return Token{TokenUnknown, 0, ""} // todo: handle heading greater than 6 levels.
    }
    return Token{TokenHeading, i, strings.Trim(line[i:], " ")}
}

func parseParagraph(lines []string) ([]string, Token) {
    value := ""
    for len(lines) > 0 {
	line := lines[0]
	if len(line) == 0 {
	    break
	}
	if len(value) > 0 {
	    value += " "
	}
	value += line
	lines = lines[1:]
    }
    return lines, Token{TokenParagraph, 0, value}
}

func parse(source string) []Token {
    lines := strings.Split(source, "\n")
    tokens := []Token{}

    for len(lines) > 0 { // todo: factor this for loop to a parseToken function.
	line := lines[0]
	token := Token{}

	if len(line) > 0 {
	    switch line[0] {
		case '#':
		    token = parseHeading(line)
		default:
		    lines, token = parseParagraph(lines)
	    }

	}

	lines = lines[1:] // fix: when paragraph eat lines we should not do it here.
	tokens = append(tokens, token)
    }
    return tokens
}

func main() {
    filePath := "./input.md"

    content, err := readFile(filePath)
    if err != nil {
	panic("failed to read file")
    }

    tokens := parse(content)
    for _, token := range tokens {
	fmt.Printf(token.toString())
    }
}
