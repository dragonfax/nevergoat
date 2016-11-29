package md

type Atom interface{
    String() string
}

type File struct {
    Parent
}

func NewFile(tr *TokenReader) *File {
    file := File{}

    for {

        token := tr.NextToken()
        switch token {
            case "":
                // end of file
                break
            case isHeader(token):
                addAtom(NewHeader(len(token), tr)
            case isAltHeader(token):
                addAtom(NewAltHeader(token))
                // because of the required K look-ahead,
                // we'll just save the atom for now, 
                // later we'll modify the generated ATL
            case isListITem(token):
                addAtom(NewList(token, tr)))
            case isBold(token):
                addAtom(NewBold(tr))
            case isItalics(token):
                addAtom(NewItalics(tr))
            case isCodeBlock(token):
                addAtom(NewCodeBlock(tr))
            case isCode(token):
                addAtom(NewCode(tr))
            // case isNumber(token):
            // case isNewline(token):
                // these aren't special in top level of a file.
            default:
                // just text, nothing special
                // TODO perhaps merge these into large blocks of text?
                addAtom(NewText(token))
        }

    }

    // TODO process file for alt headers

    return file
}

func NewHeader(rank int, tr *TokenReader) *Header {
    header := Header{make([]Atom,0,1),rank}

    for {
        token := tr.NextToken()
        switch token {
            case "", isNewLine(token):
                // end of header, or file
                return &header
            case isBold(token):
                addAtom(NewBold(tr))
            case isItalics(token):
                addAtom(NewItalics(tr))
            default:
                // no other features in headers, 
                // anything else is just a text block.
                addAtom(NewText(token))
        }
    }
}

func NewBold(tr *TokenReaader) *Bold {
    bold := Bold{make([]Atom, 0, 1)}

    for {
        token := tr.NextToken()
        switch token {
            case isBold(token):
                return &bold
            case isItalics(token):
                addAtom(NewItalics(tr))
            default:
                addAtom(NewText(token))
        }
    }
}

func NewItalics(tr *TokenReader) *Italics {
    italics := Italic{make([]Atom, 0, 1)}

    for {
        token := tr.NextToken()
        switch token {
            case isItalics(token):
                return &italics
            case isBold(token):
                addAtom(NewBold(tr))
            default:
                addAtom(NewText(token))
        }
    }
}

func NewText(token string) *Text {
    return &Text{token}
}

type Parent struct {
    children []Atom
}

func (p *Parent) addAtom(a Atom) {
    p.children = append(p.children, a)
}

type Text struct {
    content string
}

type Code struct {
    content string
}

type HorizontalRule {}

type Checkox struct {
    checked bool
}

type Header struct {
    Parent

	rank     int
}

type Bold struct {
    Parent
}

type ListITem struct {
    Parent
}

type List struct {
	list        []ListItem
	orderedList bool
}
