package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	msg := `gke-iap: a gcloud wrapper to access GKE private cluster through IAP

Usage: gke-iap [Options...]
Options:
    -p,   -project              The Google Cloud project ID to use for this invocation.
    -c,   -cluster              Name of the GKE cluster. Default: default.
    -ip,  -instance-port        The number of the GKE instance's port to connect to.
                                Default: 30443.
    -lp,  -local-port           The number of the port on which should listen for
                                connections that should be tunneled. Default: 6443.
    -v,   -version              Show version information.
    -h,   -help                 Show help message.

Author:
    Raman Antanevich <r.antanevich@gmail.com>
    <https://github.com/rantanevich>
`
	fmt.Printf(msg)
}

type Options struct {
	ProjectID    string
	ClusterName  string
	InstancePort int
	LocalPort    int
}

func ParseOptions(versionInfo string) *Options {
	opts := &Options{}

	var version bool

	flag.StringVar(&opts.ProjectID, "p", "", "Google Cloud project ID")
	flag.StringVar(&opts.ProjectID, "project", "", "Google Cloud project ID")
	flag.StringVar(&opts.ClusterName, "c", "default", "Name of the GKE cluster")
	flag.StringVar(&opts.ClusterName, "cluster", "default", "Name of the GKE cluster")
	flag.IntVar(&opts.InstancePort, "ip", 30443, "The number of the GKE instance's port to connect to")
	flag.IntVar(&opts.InstancePort, "instance-port", 30443, "The number of the GKE instance's port to connect to")
	flag.IntVar(&opts.LocalPort, "lp", 6443, "The local port on which gcloud should listen for connections that should be tunneled")
	flag.IntVar(&opts.LocalPort, "local-port", 6443, "The local port on which gcloud should listen for connections that should be tunneled")
	flag.BoolVar(&version, "v", false, "Show version information")
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.Usage = usage
	flag.Parse()

	if version {
		fmt.Println(versionInfo)
		os.Exit(0)
	}

	errMsg := opts.check()
	if len(errMsg) == 1 {
		fmt.Fprintf(os.Stderr, "Config error: %s\n", errMsg[0])
		os.Exit(1)
	} else if len(errMsg) > 1 {
		fmt.Fprintln(os.Stderr, "Config error:")
		for i, msg := range errMsg {
			fmt.Fprintf(os.Stderr, "%d. %s\n", i+1, msg)
		}
		os.Exit(1)
	}

	return opts
}

func (o *Options) check() (errMsg []string) {
	defaultProject := GetActiveProject()
	if o.ProjectID == "" && defaultProject == "" {
		errMsg = append(errMsg, "-p, -project: must be specified")
	} else if o.ProjectID == "" && defaultProject != "" {
		o.ProjectID = defaultProject
	}

	if o.ClusterName == "" {
		errMsg = append(errMsg, "-c, -cluster: cannot be empty")
	}

	if o.InstancePort > 65535 || o.InstancePort <= 0 {
		errMsg = append(errMsg, "Available INSTANCE_PORT range is 1-65535")
	}

	if o.LocalPort > 65535 || o.LocalPort < 0 {
		errMsg = append(errMsg, "Available LOCAL_PORT range is 0-65535")
	}

	return
}
