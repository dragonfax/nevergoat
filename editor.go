package main

import (
	"os"
	"os/exec"
	"sync"
	"time"
)

const ReactionTime = 200 * time.Millisecond // fast time to react to anyting

const MaximumUpdateTime = 20 * time.Second

type Cmd interface {
	Start() error
	Wait() error
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
