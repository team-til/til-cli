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
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	pb "github.com/team-til/til/server/_proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:10000"
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
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTilServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fileName := latestNote()
	fmt.Println(fileName)
	note := fileToNote(fileName)

	request := &pb.CreateNoteRequest{Note: note}
	_, err = c.CreateNote(ctx, request)
	if err != nil {
		er("Could not create note")
	}
	return nil
}

func postNote(name string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTilServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	file := findNoteByName(name)
	note := fileToNote(file)

	request := &pb.CreateNoteRequest{Note: note}
	_, err = c.CreateNote(ctx, request)
	if err != nil {
		er("Could not create note")
	}
	return nil
}

func fileToNote(fileName string) *pb.Note {
	note := &pb.Note{Filename: fileName}

	nameParts := strings.Split(fileName, ".")
	note.Name = nameParts[1]

	filePath := fmt.Sprintf("%s%s", getNotesDirectory(), fileName)
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		er("Unable to read file")
	}
	note.Note = string(fileContents)

	return note
}

func findNoteByName(name string) string {
	notesDir := getNotesDirectory()

	files, err := ioutil.ReadDir(notesDir)
	if err != nil {
		er("Unable to read notes from notes directory")
	}

	fileMatches := findByName(files, name)
	if len(fileMatches) > 1 {
		er("Found more than one file matching name:", strings.Join(fileMatches, ", "))
	} else if len(fileMatches) == 0 {
		er("Unable to find matching note")
	}
	return fileMatches[0]
}

func findByName(files []os.FileInfo, noteName string) (ret []string) {
	for _, file := range files {
		fileName := file.Name()
		namePart := strings.Split(fileName, ".")
		if namePart[1] == noteName {
			ret = append(ret, fileName)
		}
	}
	return
}

func findByTimestamp(files []os.FileInfo, ts int) (ret string) {
	for _, file := range files {
		fileName := file.Name()
		namePart := strings.Split(fileName, ".")
		s := strconv.Itoa(ts)
		if namePart[0] == s {
			ret = fileName
		}
	}
	return
}

func fileTimestamps(files []os.FileInfo) (ret []int) {
	for _, file := range files {
		fileName := file.Name()
		namePart := strings.Split(fileName, ".")
		i, err := strconv.Atoi(namePart[0])
		if err == nil {
			ret = append(ret, i)
		}
	}
	return
}

func maxIntSlice(v []int) (m int) {
	if len(v) > 0 {
		m = v[0]
	}
	for _, e := range v {
		if e < m {
			m = e
		}
	}
	return
}

func latestNote() string {
	notesDir := getNotesDirectory()

	files, err := ioutil.ReadDir(notesDir)
	if err != nil {
		er("Unable to read notes from notes directory")
	}

	timestamps := fileTimestamps(files)
	if len(timestamps) == 0 {
		er("No notes exist")
	}
	maxTimestamp := maxIntSlice(timestamps)
	return findByTimestamp(files, maxTimestamp)
}

var errPostUsageMsg = `Usage:
$ til post note-name`
