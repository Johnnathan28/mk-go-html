package main

import (
    "io"
    "os"
    "fmt"
    "strings"
)

type TokenKind int

const (
    TokenUnknown TokenKind = iota
    TokenH1 
)

func (kind TokenKind) toString() string {
    switch kind {
	case TokenUnknown: return "TokenUnknown" 
	case TokenH1:      return "TokenH1"
    }
    return ""
}

type Token struct {
    kind  TokenKind
    value string
}

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

func main() {
    filePath := "./input.md"

    content, err := readFile(filePath)
    if err != nil {
	panic("failed to read file")
    }

    tokens := []Token{}
    for _, line := range strings.Split(content, "\n") {
	token := Token{TokenUnknown, line}

	if len(line) > 0 && line[0] == '#' {
	    token.kind = TokenH1
	}

	tokens = append(tokens, token)
    }

    for _, token := range tokens {
	fmt.Printf("Token{%s, '%s'}\n", token.kind.toString(), token.value)
    }
}
