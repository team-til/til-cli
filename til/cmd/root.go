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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "til",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home := homeDir()
	viper.AddConfigPath(home)
	viper.SetConfigName(".til")

	viper.AutomaticEnv()

	viper.ReadInConfig()
}

var initRequiredMsg = `Run the init command to setup TIL:
$ til init`

func getNotesDirectory() string {
	home := homeDir()
	notesDir := viper.GetString("notesDirectory")

	if lastChar := notesDir[len(notesDir)-1:]; lastChar != "/" {
		notesDir = notesDir + "/"
	}

	notesDir = strings.Replace(notesDir, "$HOME", home, 1)

	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(notesDir, os.ModePerm); err != nil {
			er("Unable to create base notes directory")
		}
	}

	return notesDir
}
