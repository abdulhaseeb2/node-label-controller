package metadata

import "time"

const (
	FlatCarContainerLinuxLabelKey   string = "k8c.io/uses-container-linux"
	FlatCarContainerLinuxLabelValue string = "true"

	FlatCarContainerLinuxOSImageName string = "Flatcar Container Linux"

	// Requeue time
	NodeRequeueTime = 120 * time.Second
)
