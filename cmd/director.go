package cmd

type director struct {
	builder iScanBuilder
}

func newDirector(d iScanBuilder) *director {
	return &director{
		builder: d,
	}
}

func (d *director) eipBuilderScan(profiles []string) eipScan {
	d.builder.setProfiles(profiles)
	d.builder.setRegions()
	d.builder.setTargets()
	d.builder.setPorts()

	return d.builder.getEipScan()
}
