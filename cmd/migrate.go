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
	"os"

	"github.com/acheraime/certutils/backend"
	"github.com/acheraime/certutils/migrator"
	"github.com/spf13/cobra"
)

var (
	inDir              string
	outDir             string
	inCert             string
	inKey              string
	excluded           string
	destinationBackend string
	k8sCluster         string
	k8sProvider        string
	k8sNamespace       string
	projectID          string
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
		m, err := migrator.NewMigrator(backend.TLSBackendType(destinationBackend))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if inDir != "" {
			m.SetSourceDir(&inDir)
		} else {
			m.SetCertFiles(inKey, inCert)
		}

		if destinationBackend == string(backend.Backendkubernetes) {
			m.SetK8sCluster(&k8sCluster)
			m.SetBackendProvider(k8sProvider)
			m.SetProjectID(&projectID)
			m.SetK8sNamespace(&k8sNamespace)
		}

		if destinationBackend == string(backend.BackendLocal) {
			m.SetBackendDirectory(&outDir)
		}

		if err := m.Migrate(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringVarP(&inDir, "indir", "i", "", "source directory where cert files and keys are located")
	migrateCmd.Flags().StringVarP(&outDir, "outdir", "o", "", "destination directory where to move certificates. Only relevant with 'local' backend")
	migrateCmd.Flags().StringVarP(&inCert, "cert", "c", "", "path to certificate file in PEM format")
	migrateCmd.Flags().StringVarP(&inKey, "key", "k", "", "path to key file in PEM format")
	migrateCmd.Flags().StringVarP(&destinationBackend, "backend", "b", "", "certificate backend type. possible values are: local, hashivault, kubernetes")
	migrateCmd.Flags().StringVar(&k8sProvider, "k8s-provider", "", "kubernetes cloud provider, specify docker-desktop for docker-desktop")
	migrateCmd.Flags().StringVar(&k8sCluster, "k8s-cluster", "", "kubernetes cluster name, specify docker-desktop for docker-desktop")
	migrateCmd.Flags().StringVar(&k8sNamespace, "k8s-namespace", "default", "kubernetes namespace. default to 'default'")
	migrateCmd.Flags().StringVar(&projectID, "project-id", "", "cloud project ID")
}
