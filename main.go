package main

import (
	"fmt"
	"log"
)

var (
	pkgs = []string{"gcloud", "kubectl"}

	version     = "unknown"
	commit      = "unknown"
	date        = "unknown"
	versionInfo = fmt.Sprintf("Version: %s, Commit: %s, Build Time: %s", version, commit, date)
)

func main() {
	CheckPackagesInstalled(pkgs)

	gcloud := &gcloud{
		opts: ParseOptions(versionInfo),
	}

	if err := gcloud.SetupKubectl(); err != nil {
		log.Fatalln(err)
	}

	if err := gcloud.StartTunnel(); err != nil {
		log.Fatalln(err)
	}
}
