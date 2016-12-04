package main

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

type MockCmd struct {
	M sync.Mutex
}

func (mc *MockCmd) Start() error {
	return nil
}

func (mc *MockCmd) Wait() error {
	mc.M.Lock()
	mc.M.Unlock()
	return nil
}

func TestWatchEditor(t *testing.T) {

	fileUpdated := false
	updater := func() { fileUpdated = true }

	cmd := &MockCmd{}
	cmd.M.Lock()

	tempFile, _ := ioutil.TempFile("", "test")
	defer os.Remove(tempFile.Name())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		watchEditor(20*time.Millisecond, cmd, tempFile.Name(), updater)
	}()

	// to ensure timestamp differs
	time.Sleep(time.Second)

	// modify file
	tempFile.WriteString("This is a test")
	tempFile.Close()

	time.Sleep(time.Second)

	// verify updated
	if !fileUpdated {
		t.Log("failed to update on modified file")
		t.Fail()
	}

	// stop mock command
	cmd.M.Unlock()

	// wait for watchEditor to finish
	wg.Wait()
}
