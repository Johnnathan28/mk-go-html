package main

import (
    "io"
    "os"
    "fmt"
    "strings"
    "unicode"
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

func preffixLines(s string, with string) string {
    result := ""
    for _, line := range strings.Split(s, "\n") {
	if len(line) > 0 {
	    result += with + line + "\n"
	}
    }
    return result
}

func indentString(s string) string {
    return preffixLines(s, "    ")
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
    ElemListItem
    ElemOrderedList
    ElemCount
)

func (kind ElemKind) toString() string {
    switch kind {
    case ElemText:
	return "ElemText"
    case ElemBold:
	return "ElemBold"
    case ElemItalic:
	return "ElemItalic"
    case ElemBoldItalic:
	return "ElemBoldItalic"
    case ElemHeading:
	return "ElemHeading"
    case ElemParagraph:
	return "ElemParagraph"
    case ElemBlockQuote:
	return "ElemBlockQuote"
    case ElemUnknown:
	return "ElemUnknown"
    default:
	panic("Unknown element name")
    }
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
    s += fmt.Sprintf(" \"%s\"\n", e.value)
    for i, inner := range e.inner {
	s += fmt.Sprintf(" %d: [%s:\"%s\"]\n", i, inner.kind.toString(), inner.value)
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
    return lines, Element{ElemParagraph, 0, value, []Element{}}
}

func parseOrderedList(lines []string) ([]string, Element) {
    items := []Element{}
    for len(lines) > 0 {
	line := lines[0]
	if len(line) < 2 || !unicode.IsDigit(rune(line[0])) || line[1] != '.' {
	    break
	}
	line = strings.Trim(line[2:], " ")
	items = append(items, Element{ElemListItem, 0, line, []Element{}})
	lines = lines[1:]
    }
    return lines, Element{ElemOrderedList, 0, "", items}
}

func parseElement(lines []string) ([]string, Element) {
    line := lines[0]
    elem := Element{}

    if len(line) == 0 {
	lines = lines[1:]
	return lines, elem
    }

    switch c := line[0]; {
    case c == '#':
	lines, elem = parseHeading(lines)
    case c == '>':
	lines, elem = parseBlockQuote(lines)
    case unicode.IsDigit(rune(c)):
	lines, elem = parseOrderedList(lines)
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

func intoHTML(elements []Element) string {
    result := ""
    for _, elem := range elements {
	switch elem.kind {
	case ElemUnknown:
	    result += ""
	case ElemText:
	    result += elem.value
	case ElemBold:
	    panic("todo: intoHTML: ElemBold")
	case ElemItalic:
	    panic("todo: intoHTML: ElemItalic")
	case ElemBoldItalic:
	    panic("todo: intoHTML: ElemBoldItalic")
	case ElemHeading:
	    fmtString := "<h%d>%s</h%d>\n"
	    result += fmt.Sprintf(fmtString, elem.level, elem.value, elem.level)
	case ElemParagraph:
	    fmtString := "<p>%s</p>\n"
	    result += fmt.Sprintf(fmtString, elem.value)
	case ElemBlockQuote:
	    panic("todo: intoHTML: ElemBlockQuote")
	case ElemListItem:
	    fmtString := "<li>%s</li>\n"
	    result += fmt.Sprintf(fmtString, elem.value)
	case ElemOrderedList:
	    fmtString := "<ul>\n%s</ul>\n"
	    result += fmt.Sprintf(fmtString, indentString(intoHTML(elem.inner)))
	default:
	    panic("intoHTML: Unknown element kind")
	}
    }
    return result
}

func test() {
    dir := "./tests"

    files, err := os.ReadDir(dir)
    if err != nil {
	fmt.Printf("ERROR: %s\n", err)
	os.Exit(1)
    }

    for _, file := range files {
	fileName, isMarkdown := strings.CutSuffix(file.Name(), ".md")
	if !isMarkdown {
	    continue
	}

	fileNameB := ""
	for _, fileb := range files {
	    isHtml := strings.HasSuffix(fileb.Name(), ".html")
	    isEquivalent := strings.HasPrefix(fileb.Name(), fileName)

	    if isHtml && isEquivalent {
		fileNameB = fileb.Name()
	    }
	}

	if fileNameB == "" {
	    fmt.Printf("ERROR: missing equivalent HTML file for: %s\n", fileName)
	    os.Exit(1)
	}

	src, err := readFile(dir + "/" + file.Name())
	if err != nil {
	    fmt.Printf("ERROR: %s\n", err)
	    os.Exit(1)
	}

	elems := parse(src)
	html := intoHTML(elems)

	fileContentB, err := readFile(dir + "/" + fileNameB)
	if err != nil {
	    fmt.Printf("ERROR: %s\n", err)
	    os.Exit(1)
	}

	fmt.Printf("Compare '%s' with '%s'\n", file.Name(), fileNameB)
	if html != fileContentB {
	    fmt.Printf(" - Fail\n")
	    fmt.Printf("GOT:\n^%s$\nEXPECTED:\n^%s$\n", html, fileContentB)
	} else {
	    fmt.Printf(" - Success\n")
	}
    }

    os.Exit(0)
}

func main() {
    args := os.Args
    _, args = args[0], args[1:]

    if len(args) == 0 {
	filePath := "./input.md"

	content, err := readFile(filePath)
	if err != nil {
	    panic("failed to read file")
	}

	elems := parse(content)
	for _, elem := range elems {
	    fmt.Printf("%s\n", elem.toString())
	}
	os.Exit(0)
    } else if args[0] == "test" {
	test()
	os.Exit(1)
    } else {
	fmt.Printf("ERROR: Unknown command: %s\n", args[0])
	os.Exit(1)
    }
}
