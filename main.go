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
    TokenBlockQuote
    TokenCount
)

func (kind TokenKind) toString() string {
    if TokenCount > 4 {
	panic("missing name of new token added")
    }
    switch kind {
	case TokenHeading: return "TokenHeading"
	case TokenParagraph: return "TokenParagraph"
	case TokenBlockQuote: return "TokenBlockQuote"
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

func parseHeading(lines []string) ([]string, Token) {
    line := lines[0]
    i := 0
    for ; i < len(line); i += 1 {
	if line[i] != '#' {
	    break
	}
    }
    lines = lines[1:]
    if i > 6 {
	// note: Not so sure about how handle more than 6 hashtags for titles
	// it should be handled as a paragraph? or a error? should it be ignored?
	return lines, Token{TokenUnknown, 0, ""}
    }
    return lines, Token{TokenHeading, i, strings.Trim(line[i:], " ")}
}

func parseBlockQuote(lines []string) ([]string, Token) {
    value := ""
    newLine := false
    for len(lines) > 0 && len(lines[0]) > 0 && lines[0][0] == '>' {
	lineValue := strings.Trim(lines[0][1:], " ")
	if len(lineValue) > 0 {
	    if len(value) > 0 && !newLine {
		value += " "
	    }
	    newLine = false
	} else {
	    value += "  " // todo: this should be a newline.
	    newLine = true
	}
	value += lineValue
	lines = lines[1:]
    }
    return lines[1:], Token{TokenBlockQuote, 0, value}
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

func parseToken(lines []string) ([]string, Token) {
    line := lines[0]
    token := Token{}

    if len(line) == 0 {
	lines = lines[1:]
	return lines, token
    }

    switch line[0] {
    case '#':
	lines, token = parseHeading(lines)
    case '>':
	lines, token = parseBlockQuote(lines)
    default:
	lines, token = parseParagraph(lines)
    }
    return lines, token
}

func parse(source string) []Token {
    lines := strings.Split(source, "\n")
    tokens := []Token{}

    for len(lines) > 0 {
	token := Token{}
	lines, token = parseToken(lines)
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
