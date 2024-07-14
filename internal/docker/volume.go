package docker

import (
	"encoding/json"
)

func (volume *Volume) Load() error {
	output, err := DockerVolumeExec("inspect", "--format", "json", volume.Name)
	if err != nil || !json.Valid([]byte(output)) {
		return err
	}
	var jsonOutput Volume
	json.Unmarshal([]byte(output), &jsonOutput)
	return nil
}

func (volume *Volume) Create() error {
	args := []string{"create"}
	if volume.Driver == "" {
		volume.Driver = "local"
	}
	args = append(args, "--driver", volume.Driver)
	volume.Labels["creator"] = "dockonaut"
	for key, value := range volume.Labels {
		args = append(args, "--label", key+"="+value)
	}
	for key, value := range volume.Options {
		args = append(args, "--opt", key+"="+value)
	}
	args = append(args, volume.Name)
	_, err := DockerVolumeExec(args...)
	return err
}

func (volume *Volume) Drop() error {
	_, err := DockerVolumeExec("rm", volume.Name)
	return err
}
