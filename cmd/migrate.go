/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	"github.com/acheraime/certutils/backend"
	"github.com/acheraime/certutils/migrator"
	"github.com/spf13/cobra"
)

var (
	sourceDir          string
	destinationDir     string
	destinationBackend string
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrator.NewMigrator(backend.TLSBackendType(destinationBackend)); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")
	migrateCmd.Flags().StringVar(&sourceDir, "in", "", "source directory where cert files and keys are located")
	migrateCmd.Flags().StringVar(&destinationDir, "out", "", "destination directory where to move certificates. Only relevant with local backend")
	migrateCmd.Flags().StringVarP(&destinationBackend, "backend", "b", "", "certificate backend type. possible values are: local, hashivault, kubernetes")
}
