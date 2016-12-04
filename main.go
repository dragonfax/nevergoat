package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"time"

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

func getTempFile() string {
	tempFile, err := ioutil.TempFile("", "gote")
	if err != nil {
		panic("failed to get a temp file: " + err.Error())
	}
	err = tempFile.Close()
	if err != nil {
		panic("failed to close temp file")
	}

	return tempFile.Name()
}

func startEditor(tempFileName string) *exec.Cmd {
	cmd := exec.Command("vim", tempFileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic("Editor returned error: " + err.Error())
	}

	return cmd
}

const ReactionTime = 200 * time.Millisecond // fast time to react to anyting

type Cmd interface {
	Start() error
	Wait() error
}

func watchEditor(updateTime time.Duration, cmd Cmd, tempFileName string, updater func()) {

	wg := sync.WaitGroup{}

	editorDone := false

	wg.Add(1)
	go func() {
		defer wg.Done()

		stat, _ := os.Stat(tempFileName)
		initialTime := stat.ModTime()

		// cross platform fsnotify is spotty for now.

		lastUpdateCheck := time.Now()

		for !editorDone {

			// Fastest time to respone to anything.
			time.Sleep(ReactionTime)

			// check update time.
			thisUpdateCheck := time.Now()
			if thisUpdateCheck.Sub(lastUpdateCheck) > updateTime {
				if stat, _ = os.Stat(tempFileName); initialTime != stat.ModTime() {
					updater()
				}
			}
			lastUpdateCheck = thisUpdateCheck
		}

	}()

	cmd.Wait()

	wg.Wait()

	// one last update
	updater()
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

const MaximumUpdateTime = 20 * time.Second

func main() {

	readSettings()

	enClient := connect()

	tempFileName := getTempFile()
	defer func() {
		if err := os.Remove(tempFileName); err != nil {
			log.Println("failed to remove the temp file ", err)
		}
	}()

	cmd := startEditor(tempFileName)

	watchEditor(MaximumUpdateTime, cmd, tempFileName, func() { updateChanges(enClient, tempFileName) })
}
