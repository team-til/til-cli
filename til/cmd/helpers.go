package cmd

import (
	"bytes"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
)

func er(msgs ...string) {
	var buffer bytes.Buffer

	buffer.WriteString("Error: ")
	buffer.WriteString(msgs[0])
	buffer.WriteString("\n")

	for i, msg := range msgs {
		if i != 0 {
			buffer.WriteString("\n")
			buffer.WriteString(msg)
			buffer.WriteString("\n")
		}
	}

	fmt.Println(buffer.String())
	os.Exit(1)
}

func homeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}
