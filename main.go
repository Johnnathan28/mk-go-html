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

type ElemKind int

const (
    ElemUnknown ElemKind = iota
    ElemText
    ElemBold
    ElemItalic
    ElemBoldItalic
    ElemHeading
    ElemParagraph
    ElemBlockQuote
    ElemCount
)

func (kind ElemKind) toString() string {
    if ElemCount > 8 {
	panic("missing name of new token added")
    }
    switch kind {
	case ElemText: return "ElemText"
	case ElemBold: return "ElemBold"
	case ElemItalic: return "ElemItalic"
	case ElemBoldItalic: return "ElemBoldItalic"
	case ElemHeading: return "ElemHeading"
	case ElemParagraph: return "ElemParagraph"
	case ElemBlockQuote: return "ElemBlockQuote"
    }
    return "ElemUnknown" 
}

type Element struct {
    kind  ElemKind
    level int
    value string
    inner []Element
}

func newElement() Element {
    return Element{ElemUnknown, 0, "", []Element{}}
}

func (e *Element) toString() string {
    s := fmt.Sprintf("Element(%s:%d)\n", e.kind.toString(), e.level)
    s += fmt.Sprintf(" '%s'\n", e.value)
    for i, inner := range e.inner {
	s += fmt.Sprintf(" %d: [%s:'%s']\n", i, inner.kind.toString(), inner.value)
    }
    return s
}

func parseHeading(lines []string) ([]string, Element) {
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
	return lines, newElement()
    }
    return lines, Element{ElemHeading, i, strings.Trim(line[i:], " "), []Element{}}
}

func parseBlockQuote(lines []string) ([]string, Element) {
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
    return lines[1:], Element{ElemBlockQuote, 0, value, []Element{}}
}

func parseParagraph(lines []string) ([]string, Element) {
    value := ""

    for len(lines) > 0 && len(lines[0]) > 0 {
	if len(value) > 0 {
	    value += " "
	}
	value += lines[0]
	lines = lines[1:]
    }

    bold := false
    italic := false
    boldItalic := false

    inner := []Element{}

    tmp := ""
    buff := value

    for len(buff) > 0 {
	if strings.HasPrefix(buff, "***") || strings.HasPrefix(buff, "___") {
	    buff = buff[3:]
	    if boldItalic {
		inner = append(inner, Element{ElemBoldItalic, 0, tmp, []Element{}})
	    } else {
		inner = append(inner, Element{ElemText, 0, tmp, []Element{}})
	    }
	    boldItalic = !boldItalic
	    tmp = ""
	} else if strings.HasPrefix(buff, "**") || strings.HasPrefix(buff, "__") {
	    buff = buff[2:]
	    if bold {
		inner = append(inner, Element{ElemBold, 0, tmp, []Element{}})
	    } else {
		inner = append(inner, Element{ElemText, 0, tmp, []Element{}})
	    }
	    bold = !bold
	    tmp = ""
	} else if strings.HasPrefix(buff, "*") || strings.HasPrefix(buff, "_") {
	    buff = buff[1:]
	    if italic {
		inner = append(inner, Element{ElemItalic, 0, tmp, []Element{}})
	    } else {
		inner = append(inner, Element{ElemText, 0, tmp, []Element{}})
	    }
	    italic = !italic
	    tmp = ""
	}

	tmp += buff[:1]
	buff = buff[1:]
    }

    inner = append(inner, Element{ElemText, 0, tmp, []Element{}})
    return lines, Element{ElemParagraph, 0, value, inner}
}

func parseElement(lines []string) ([]string, Element) {
    line := lines[0]
    elem := Element{}

    if len(line) == 0 {
	lines = lines[1:]
	return lines, elem
    }

    switch line[0] {
    case '#':
	lines, elem = parseHeading(lines)
    case '>':
	lines, elem = parseBlockQuote(lines)
    default:
	lines, elem = parseParagraph(lines)
    }
    return lines, elem
}

func parse(source string) []Element {
    lines := strings.Split(source, "\n")
    elems := []Element{}

    for len(lines) > 0 {
	elem := Element{}
	lines, elem = parseElement(lines)
	elems = append(elems, elem)
    }
    return elems
}

func main() {
    filePath := "./input.md"

    content, err := readFile(filePath)
    if err != nil {
	panic("failed to read file")
    }

    elems := parse(content)
    for _, elem := range elems {
	fmt.Printf("%s\n", elem.toString())
    }
}
