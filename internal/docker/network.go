package docker

import (
	"encoding/json"
	"errors"
	"strings"
)

func (network *Network) Load() error {
	if network.Name == "" {
		return errors.New("network name is empty")
	}
	output, err := DockerNetworkExec("inspect", "--format", "json", network.Name)
	if err != nil || !json.Valid([]byte(output)) {
		return errors.Join(errors.New("docker exec error"), err)
	}

	// output is in an array
	output = strings.TrimSpace(output)
	output = strings.TrimPrefix(output, "[")
	output = strings.TrimSuffix(output, "]")

	if err := json.Unmarshal([]byte(output), &network); err != nil {
		return errors.Join(errors.New(output[len(output)-10:]), err)
	}
	return nil
}

func (network *Network) Create() error {
	args := []string{"create"}
	args = append(args, "--driver", network.Driver)
	if network.Labels == nil {
		network.Labels = map[string]string{}
	}
	network.Labels["creator"] = "dockonaut"
	for key, value := range network.Labels {
		args = append(args, "--label", key+"="+value)
	}
	for key, value := range network.Options {
		args = append(args, "--opt", key+"="+value)
	}
	if network.Attachable {
		args = append(args, "--attachable")
	}
	if network.EnableIPv6 {
		args = append(args, "--ipv6")
	}
	if network.Internal {
		args = append(args, "--internal")
	}
	for _, ipam := range network.IPAM.Config {
		if ipam.Subnet != "" {
			args = append(args, "--subnet="+ipam.Subnet)
		}
		if ipam.IPRange != "" {
			args = append(args, "--ip-range="+ipam.IPRange)
		}
		if ipam.Gateway != "" {
			args = append(args, "--gateway="+ipam.Gateway)
		}
	}

	args = append(args, network.Name)
	_, err := DockerNetworkExec(args...)
	return err
}

func (network *Network) Drop() error {
	_, err := DockerNetworkExec("rm", network.Name)
	return err
}
