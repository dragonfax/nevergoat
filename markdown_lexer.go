package main

// Simple markdown <-> enml translation

import (
	"os"
)


type TokenReader struct {
    runeReader strings.Reader
}

func NewTokenReader(s string) *TokenReader {
    return &TokenReader{strings.NewReader(s)}
}

func (tr *TokenReader) NextToken() string {

    r, _, err := tr.runeReader.ReadRune()

    if err != nil {
        if err == io.EOF {
            reeturn ""
        } else {
            log.Fatalf("Failure in tokenzing Markdown stream: %v", err)
        }
    }

    switch r {
        case '\n', '.', ' ':
            return string(r)
        case '=', '-', '#', '*', '_':
            return continueToRead(tr.runeReader,r, r)
        case '0','1','2','3','4','5','6','7','8','9'
            return continueToRead(tr.runReader,r,"0123456789")
        default:
            // gather text chars into larger groups.
            return continueToReadOther(tr.runeReader, r, "\n=-#*_ ")

    }

}

func continueToRead(reader strings.Reader, firstRune rune, toRead rune) string {

}
