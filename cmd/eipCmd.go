/*
Copyright © 2021 Stefano Pirrello spirrello@gmail.com

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
	"strings"

	"github.com/spf13/cobra"
)

// eipCmd represents the eip command
var eipCmd = &cobra.Command{
	Use:   "eip",
	Short: "perform a scan against AWS Elastic IPs",
	Long:  `perform a scan against AWS Elastic IPs.`,
	Run: func(cmd *cobra.Command, args []string) {

		var awsConfig, tempProfiles string
		var profiles []string

		tempProfiles, _ = cmd.Flags().GetString("profiles")

		profiles = strings.Split(tempProfiles, ",")

		if profiles[0] == "" {
			awsConfig, _ = cmd.Flags().GetString("awsConfig")
			profiles = loadAWSConfigProfiles(awsConfig)
		}

		fmt.Printf("profiles: %s", profiles)

		builder := getScanBuilder("eip")
		director := newDirector(builder)
		e := director.eipBuilderScan(profiles)

		fmt.Printf("targets: %s\n", e.scanTargets)

	},
}

func init() {
	rootCmd.AddCommand(eipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	eipCmd.PersistentFlags().String("awsConfig", "~/.aws/config", "AWS config file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	eipCmd.Flags().String("profiles", "", "which account to scan")
}
