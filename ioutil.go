package main

import "io/ioutil"

const tempTemplate = `title: Untitled
---
`

func getTempFile() string {
	tempFile, err := ioutil.TempFile("", "gote")
	if err != nil {
		panic("failed to get a temp file: " + err.Error())
	}

	tempFile.WriteString(tempTemplate)

	err = tempFile.Close()
	if err != nil {
		panic("failed to close temp file")
	}

	return tempFile.Name()
}
