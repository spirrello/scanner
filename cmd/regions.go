package cmd

type AWSRegions struct {
	Regions []struct {
		Endpoint    string `json:"Endpoint"`
		RegionName  string `json:"RegionName"`
		OptInStatus string `json:"OptInStatus"`
	} `json:"Regions"`
}
