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
	"fmt"
	"os"

	"github.com/acheraime/certutils/configurator"
	"github.com/acheraime/certutils/utils"
	"github.com/spf13/cobra"
)

var (
	cfgInFile string
	cfgOutDir string
)

// confgenCmd represents the confgen command
var confgenCmd = &cobra.Command{
	Use:   "confgen",
	Short: "Use confgen command to generate configuration for nginx and traefik proxies",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgOutDir == "" {
			cfgOutDir = utils.UserDesktop()
		}
		p, err := configurator.NewConfigurator(configurator.Options{
			In:       cfgInFile,
			Out:      cfgOutDir,
			Type:     configurator.RedirectConfig,
			Platform: configurator.TraefikPlatform,
			InType:   configurator.CSVInput,
			DryRun:   true,
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		//p.DryRun("/Users/acheraime/Desktop/testout")
		// if err := p.Generate(); err != nil {
		// 	fmt.Println("unable to generate configuration. " + err.Error())
		// 	os.Exit(1)
		// }
		if err = p.Generate(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("configuration generated successfully")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(confgenCmd)

	confgenCmd.Flags().StringVarP(&cfgInFile, "input", "i", "", "location of CSV file where to read redirects paths from")
	confgenCmd.Flags().StringVarP(&cfgOutDir, "outdir", "o", "", "Directory where to dump generated configuration file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// confgenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// confgenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
