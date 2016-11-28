package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/dragonfax/evernote-sdk-go/notestore"
	"github.com/dragonfax/evernote-sdk-go/types"
)

func strP(s string) *string {
	return &s
}

func main() {
	notestoreUrl := os.Getenv("EVERNOTE_NOTESTORE")

	trans, err := thrift.NewTHttpPostClient(notestoreUrl)
	if err != nil {
		panic(err)
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := notestore.NewNoteStoreClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		panic(err)
	}

	authenticationToken := os.Getenv("EVERNOTE_TOKEN")

	tempFile, err := ioutil.TempFile("", "gote")
	if err != nil {
		panic("failed to get a temp file: " + err.Error())
	}
	tempFile.Close()
	os.Remove(tempFile.Name())
	defer os.Remove(tempFile.Name())
	cmd := exec.Command("vim", tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic("Editor returned error: " + err.Error())
	}

	fileContent, err := ioutil.ReadFile(tempFile.Name())

	note := types.NewNote()
	note.Title = strP("Test Note")
	note.Content = strP(
		`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">
<en-note>  
	<pre> 
` + string(fileContent) + `
	</pre>
</en-note>
`)

	updatedNote, err := client.CreateNote(authenticationToken, note)
	if err != nil {
		panic(err)
	}

	fmt.Println(updatedNote)

}
