package docker

import (
	"encoding/json"
)

func (volume *Volume) Load() error {
	output, err := DockerVolumeExec("inspect", "--format", "json", volume.Name)
	if err != nil || !json.Valid([]byte(output)) {
		return err
	}
	var jsonOutput []map[string]any
	json.Unmarshal([]byte(output), &jsonOutput)
	volume.Driver= jsonOutput[0]["Driver"].(string)
	if jsonOutput[0]["Labels"] != nil {
		for i, v := range jsonOutput[0]["Labels"].(map[string]interface{}) {
			volume.Labels = append(volume.Labels, i+"="+v.(string))
		}
	}
	if jsonOutput[0]["Options"] != nil {
		for i, v := range jsonOutput[0]["Options"].(map[string]interface{}) {
			volume.Options = append(volume.Options, i+"="+v.(string))
		}
	}
	return nil
}

func (volume *Volume) Create() error {
	args := []string{"create"}
	if volume.Driver == "" {
		volume.Driver = "local"
	}
	args = append(args, "--driver", volume.Driver)
	volume.Labels = append(volume.Labels, "creator=dockonaut")
	for _, label := range volume.Labels {
		args = append(args, "--label", label)
	}
	for _, option := range volume.Options {
		args = append(args, "--opt", option)
	}
	args = append(args, volume.Name)
	_, err := DockerVolumeExec(args...)
	return err
}

func (volume *Volume) Drop() error {
	_, err := DockerVolumeExec("rm", volume.Name)
	return err
}
