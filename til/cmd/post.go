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

	"github.com/spf13/cobra"
)

var postLatest bool

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post a note to Team TIL",
	Long: `Post a note to the Team TIL server. 

Usage example:

$ til post note-name

Tip - As a shortcut you can automatically post your latest (most recently created) note with --latest or -l
`,
	Run: post,
}

func init() {
	rootCmd.AddCommand(postCmd)

	postCmd.Flags().BoolVarP(&postLatest, "latest", "l", false, "Post your most recent created note")
}

func post(cmd *cobra.Command, args []string) {
	if postLatest {
		if err := postLatestNote(); err != nil {
			er("Unable to post note")
		}
		fmt.Println("Posted note! You can see your posted note at: https://company.til.team/notes/asdf1234")
	} else if len(args) != 1 {
		er("You must provide a single note name.", errPostUsageMsg)
	} else {
		if err := postNote(args[0]); err != nil {
			er("Unable to post note")
		}
		fmt.Println("Posted note! You can see your posted note at: https://company.til.team/notes/asdf1234")
	}
}

func postLatestNote() error {
	return nil
}

func postNote(name string) error {
	return nil
}

var errPostUsageMsg = `Usage:
$ til post note-name`
