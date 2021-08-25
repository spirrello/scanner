package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/bigkevmcd/go-configparser"
)

//loadAWSConfigProfiles parses the aws config to get all config profiles
func loadAWSConfigProfiles(awsConfig string) []string {
	fmt.Printf("our awsConfig: %s\n", awsConfig)

	var profiles []string

	p, err := configparser.NewConfigParserFromFile(awsConfig)
	if err != nil {
		log.Fatal(err)
	}
	p.RemoveSection("default")
	profiles = p.Sections()

	return cleanseAWSProfiles(profiles)
}

func cleanseAWSProfiles(profiles []string) []string {

	for i, v := range profiles {
		profiles[i] = strings.ReplaceAll(v, "profile", "")
	}

	return profiles
}
