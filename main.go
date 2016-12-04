package main

import (
	"log"
	"os"
)

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
