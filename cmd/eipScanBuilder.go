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
	ports       string
	profiles    []string
	regions     []string
}

type ProfileAddresses struct {
	AWSAccount string `json:"AWSAccount"`
	Addresses  []struct {
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

func (e *eipScanBuilder) setPorts(ports string) {
	e.ports = ports
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

//setTargets will fetch all EIPs from each region/profile
func (e *eipScanBuilder) setTargets() {
	//no go-staticcheck
	for _, v := range e.profiles {
		for _, p := range e.regions {

			awsCommand := "aws ec2 describe-addresses --profile " + v + " --region " + p + " --no-paginate"
			out, err := exec.Command("bash", "-c", awsCommand).Output()

			if err != nil {
				log.Fatal(err)
			}

			var awsAddresses ProfileAddresses
			err = json.Unmarshal([]byte(out), &awsAddresses)
			if err != nil {
				fmt.Printf("error parsing aws ec2 describe-addresses - %s", err)
			}

			//assign profile to identify which account the EIPs belong to
			awsAddresses.AWSAccount = v

			if len(awsAddresses.Addresses) > 0 {
				e.scanTargets = append(e.scanTargets, awsAddresses)
			}

		}
	}
}

func (e *eipScanBuilder) setProfiles(profiles []string) {
	e.profiles = profiles
}

//runScan executes an nmap scan
func (e *eipScanBuilder) runScan(conf string) {

	for _, v := range e.scanTargets {
		for _, a := range v.Addresses {
			nmapScan := conf + " " + e.ports + " " + a.PublicIP
			out, err := exec.Command("bash", "-c", nmapScan).Output()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("###############################################")
			fmt.Printf("Account: %s\n", v.AWSAccount)
			fmt.Printf("IP: %s\n", a.PublicIP)
			if strings.Contains(string(out), "open") {
				fmt.Println(string(out))
			} else {
				fmt.Println("result: all ports closed")
			}

		}
	}
}

func (e *eipScanBuilder) getEipScan() eipScan {
	return eipScan{
		scanTargets: e.scanTargets,
		ports:       e.ports,
		profiles:    e.profiles,
	}
}
