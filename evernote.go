package main

import (
	"io/ioutil"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/dragonfax/evernote-sdk-go/notestore"
	"github.com/dragonfax/evernote-sdk-go/types"
)

func strP(s string) *string {
	return &s
}

func connect() *notestore.NoteStoreClient {

	trans, err := thrift.NewTHttpPostClient(Settings.Notestore)
	if err != nil {
		panic(err)
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := notestore.NewNoteStoreClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		panic(err)
	}

	return client
}

func updateChanges(enClient *notestore.NoteStoreClient, tempFileName string) *types.Note {

	fileContent, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		panic("failed to read temp file")
	}

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

	updatedNote, err := enClient.CreateNote(Settings.Token, note)
	if err != nil {
		panic(err)
	}

	return updatedNote
}
