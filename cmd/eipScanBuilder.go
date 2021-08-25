package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type eipScanBuilder struct {
	scanTargets []ProfileAddresses
	ports       []string
	profiles    []string
	regions     []string
}

type ProfileAddresses struct {
	Addresses []struct {
		PublicIP                string `json:"PublicIp"`
		AllocationID            string `json:"AllocationId"`
		AssociationID           string `json:"AssociationId"`
		Domain                  string `json:"Domain"`
		NetworkInterfaceID      string `json:"NetworkInterfaceId"`
		NetworkInterfaceOwnerID string `json:"NetworkInterfaceOwnerId"`
		PrivateIPAddress        string `json:"PrivateIpAddress"`
		Tags                    []struct {
			Key   string `json:"Key"`
			Value string `json:"Value"`
		} `json:"Tags"`
		PublicIpv4Pool     string `json:"PublicIpv4Pool"`
		NetworkBorderGroup string `json:"NetworkBorderGroup"`
	} `json:"Addresses"`
}

func newEipScanBuilder() *eipScanBuilder {
	return &eipScanBuilder{}
}

func (e *eipScanBuilder) setPorts() {

	e.ports = strings.Split("22,80,443,9735", ",")
}

func (e *eipScanBuilder) setRegions() {

	out, err := exec.Command("bash", "-c", "aws ec2 describe-regions").Output()
	if err != nil {
		log.Fatal(err)
	}
	var awsRegions AWSRegions
	err = json.Unmarshal([]byte(out), &awsRegions)
	if err != nil {
		fmt.Printf("error parsing aws ec2 describe-regions - %s", err)
	}

	for _, v := range awsRegions.Regions {
		e.regions = append(e.regions, v.RegionName)
	}
}

func (e *eipScanBuilder) setTargets() {

	//no go-staticcheck
	out, err := exec.Command("bash", "-c", "aws ec2 describe-addresses --profile dev --region us-east-1 --no-paginate").Output()

	if err != nil {
		log.Fatal(err)
	}
	var awsAddresses ProfileAddresses
	err = json.Unmarshal([]byte(out), &awsAddresses)
	if err != nil {
		fmt.Printf("error parsing aws ec2 describe-addresses - %s", err)
	}
	e.scanTargets = append(e.scanTargets, awsAddresses)
}

func (e *eipScanBuilder) setProfiles(profiles []string) {

	e.profiles = profiles
}

func (e *eipScanBuilder) getEipScan() eipScan {
	return eipScan{
		scanTargets: e.scanTargets,
		ports:       e.ports,
		profiles:    e.profiles,
	}
}
