package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/knieriem/markdown"
)

const enmlPrefix = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">
<en-note>
`
const enmlSuffix = `
</en-note>`

func MDToENML(md string) string {
	p := markdown.NewParser(&markdown.Extensions{Smart: true})
	w := &bytes.Buffer{}
	p.Markdown(strings.NewReader(md), markdown.ToHTML(w))
	return enmlPrefix + w.String() + enmlSuffix
}

// Simple SAX parsing of the XML

type List struct {
	ordered bool
	num     int
}

type ListTrack []*List

func (lt *ListTrack) addOrderedList() {
	*lt = append(*lt, &List{true, 0})
}

func (lt *ListTrack) addUnorderedList() {
	*lt = append(*lt, &List{false, 0})
}

func (lt *ListTrack) addListItem(buf *bytes.Buffer, token xml.Token) {
	listLen := len(*lt)

	// indent first
	for i := 1; i < listLen; i++ {
		buf.WriteString("  ")
	}

	// list item
	list := lt.last()
	if list.ordered {
		buf.WriteString(fmt.Sprintf("%d. ", list.num))
		list.num++
	} else {
		buf.WriteString("* ")
	}
}

func (lt *ListTrack) popList() {
	*lt = (*lt)[0 : len(*lt)-1]
}

func (lt *ListTrack) last() *List {
	return (*lt)[len(*lt)-1]
}

func ENMLToMD(enml string) string {

	buf := bytes.NewBuffer([]byte(""))

	list := make(ListTrack, 0, 0)

	decoder := xml.NewDecoder(strings.NewReader(enml))

	for token, err := decoder.Token(); err != io.EOF; token, err = decoder.Token() {
		if err != nil {
			log.Fatalf("some sort of error while parsing xml %v", err)
		}
		switch token := token.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "h1":
				buf.WriteString("# ")
			case "h2":
				buf.WriteString("## ")
			case "h3":
				buf.WriteString("### ")
			case "h4":
				buf.WriteString("#### ")
			case "h5":
				buf.WriteString("##### ")
			case "h6":
				buf.WriteString("###### ")
			case "b":
				buf.WriteString("**")
			case "em", "i":
				buf.WriteString("_")
			case "ul":
				list.addUnorderedList()
			case "ol":
				list.addOrderedList()
			case "li":
				list.addListItem(buf, token)
			case "br", "p", "en-note", "div", "span":
			default:
				log.Fatalf("unknown start tag %v", token.Name.Local)
			}
		case xml.EndElement:
			switch token.Name.Local {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				buf.WriteString("\n")
			case "b":
				buf.WriteString("**")
			case "em", "i":
				buf.WriteString("_")
			case "li":
				buf.WriteString("\n")
			case "ul", "ol":
				list.popList()
			case "p", "br", "span":
				buf.WriteString("\n")
			case "en-note", "div":
			default:
				log.Fatalf("unknown tag %v", token.Name.Local)
			}
		case xml.CharData:
			buf.Write([]byte(token))
		case xml.Comment, xml.ProcInst, xml.Directive:
		default:
			log.Fatalf("unknown xml token type %T", token)
		}
	}

	return buf.String()
}
