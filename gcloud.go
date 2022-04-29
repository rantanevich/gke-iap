package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type gcloud struct {
	opts *Options
}

type gkeInstance struct {
	Name string
	Zone string
}

func (g *gcloud) SetupKubectl() error {
	args := []string{
		"container",
		"clusters",
		"get-credentials",
		g.opts.ClusterName,
		"--internal-ip",
		"--project", g.opts.ProjectID,
	}

	location, err := g.getGKELocation()
	if err != nil {
		return err
	}

	if len(strings.Split(location, "-")) == 3 {
		args = append(args, "--zone", location)
	} else {
		args = append(args, "--region", location)
	}

	msg := fmt.Sprintf("Setting 'gke_%s_%s_%s' context in kubeconfig\n", g.opts.ProjectID, location, g.opts.ClusterName)
	_, err = Exec(args, msg)
	if err != nil {
		return err
	}

	return nil
}

func (g *gcloud) StartTunnel() error {
	instances, err := g.getGKEInstances()
	if err != nil {
		return err
	}

	instance := GetRandomGKEInstance(instances)

	args := []string{
		"compute",
		"start-iap-tunnel",
		"--project", g.opts.ProjectID,
		"--zone", instance.Zone,
		"--local-host-port", fmt.Sprintf("127.0.0.1:%d", g.opts.LocalPort),
		instance.Name, strconv.Itoa(g.opts.InstancePort),
	}

	msg := fmt.Sprintf("Listening on port [%d]\n", g.opts.LocalPort)
	_, err = Exec(args, msg)
	if err != nil {
		return err
	}

	return nil
}

func (g *gcloud) getGKELocation() (string, error) {
	args := []string{
		"container",
		"clusters",
		"list",
		"--project", g.opts.ProjectID,
		"--filter", fmt.Sprintf("name=%s", g.opts.ClusterName),
		"--format", "value(location)",
	}

	msg := fmt.Sprintf("Fetching '%s' GKE cluster location in '%s'\n", g.opts.ClusterName, g.opts.ProjectID)
	stdout, err := Exec(args, msg)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout), nil
}

func (g *gcloud) getGKEInstances() ([]gkeInstance, error) {
	args := []string{
		"compute",
		"instances",
		"list",
		"--project", g.opts.ProjectID,
		"--filter", fmt.Sprintf("name~^gke-%s-", g.opts.ClusterName),
		"--format", "csv[no-heading](name,zone)",
	}

	msg := fmt.Sprintf("Fetching instances of '%s' GKE cluster in '%s'\n", g.opts.ClusterName, g.opts.ProjectID)
	stdout, err := Exec(args, msg)
	if err != nil {
		return []gkeInstance{}, err
	}

	var instances []gkeInstance

	rows := strings.TrimSpace(stdout)

	for _, row := range strings.Split(rows, "\n") {
		data := strings.Split(row, ",")
		instances = append(instances, gkeInstance{
			Name: data[0],
			Zone: data[1],
		})
	}

	return instances, nil
}

func GetActiveProject() string {
	args := []string{
		"config",
		"get",
		"project",
	}
	stdout, _ := Exec(args, "")
	return strings.TrimSpace(stdout)
}

func GetRandomGKEInstance(instances []gkeInstance) gkeInstance {
	rand.Seed(time.Now().UnixNano())
	return instances[rand.Intn(len(instances))]
}
