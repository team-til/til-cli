// Copyright © 2018 Nathan Owsiany <nowsiany@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create [note-name]",
	Short: "Create and open a new TIL markdown note",
	Long: `Create generates a new TIL markdown note in the configured TIL notes directory
and opens it in your desired editor. You can configure your notes directory and
desired editor by editing the generated ~/.til.yaml config file.

Usage example:

$ til create note-name

The example above generates a new $NOTES_DIR/20180519142015.note-name.md
markdown file.`,
	Run: create,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		er("You must provide a note name.", errCreateUsageMsg)
	} else if len(args) > 1 {
		er("please provide only one argument", errCreateUsageMsg)
	} else {
		fmt.Println("Creating note...")
		path, err := createNote(args[0])
		if err != nil {
			er("Unable to create note", err.Error())
		}

		fmt.Printf("Note created: %s\n", path)
		openNoteInEditor(path)
	}
}

func createNote(noteName string) (string, error) {
	notesDir := getNotesDirectory()

	notePath := fmt.Sprintf("%s%s.%s.md", notesDir, currTimestamp(), noteName)

	file, err := os.OpenFile(notePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
		er("Unable to create file")
	}
	defer file.Close()

	return notePath, nil
}

func currTimestamp() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func openNoteInEditor(notePath string) {
	openNotesOnCreate := viper.GetBool("openNotesOnCreate")
	if openNotesOnCreate {
		fmt.Println("Opening created note...")

		editor := viper.GetString("noteEditor")

		cmd := exec.Command(editor, notePath)

		err := cmd.Run()
		if err != nil {
			er("Unable to open note")
		}
	}
}

var errCreateUsageMsg = `Usage:
$ til create note-name`
