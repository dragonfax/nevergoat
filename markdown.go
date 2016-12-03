package main

import "github.com/knieriem/markdown"
import "bytes"
import "strings"

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
