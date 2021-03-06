package cmd

type iScanBuilder interface {
	setTargets()
	setPorts(ports string)
	setProfiles(profiles []string)
	setRegions()
	getEipScan() eipScan
	runScan(conf string)
}

func getScanBuilder(scanType string) iScanBuilder {
	if scanType == "eip" {
		return &eipScanBuilder{}
	}

	// return nil, fmt.Errorf("we need a valid scan type, eip or elb")
	return nil
}
