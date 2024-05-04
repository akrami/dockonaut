package docker

import "encoding/json"

func (network *Network) Load() error {
	output, err := DockerNetworkExec("inspect", "--format", "json", network.Name)
	if err != nil || !json.Valid([]byte(output)) {
		return err
	}
	var jsonOutput []map[string]any
	json.Unmarshal([]byte(output), &jsonOutput)
	network.Driver = jsonOutput[0]["Driver"].(string)
	network.IPV6 = jsonOutput[0]["EnableIPv6"].(bool)
	network.Internal = jsonOutput[0]["Internal"].(bool)
	network.Attachable = jsonOutput[0]["Attachable"].(bool)
	network.SubNet = jsonOutput[0]["IPAM"].(map[string]interface{})["Config"].([]interface{})[0].(map[string]interface{})["Subnet"].(string)
	network.IPRange = jsonOutput[0]["IPAM"].(map[string]interface{})["Config"].([]interface{})[0].(map[string]interface{})["IPRange"].(string)
	network.Gateway = jsonOutput[0]["IPAM"].(map[string]interface{})["Config"].([]interface{})[0].(map[string]interface{})["Gateway"].(string)
	if jsonOutput[0]["Labels"] != nil {
		for i, v := range jsonOutput[0]["Labels"].(map[string]interface{}) {
			network.Labels = append(network.Labels, i+"="+v.(string))
		}
	}
	if jsonOutput[0]["Options"] != nil {
		for i, v := range jsonOutput[0]["Options"].(map[string]interface{}) {
			network.Options = append(network.Options, i+"="+v.(string))
		}
	}
	return nil
}

func (network *Network) Create() error {
	args := []string{"create"}
	args = append(args, "--driver", network.Driver)
	network.Labels = append(network.Labels, "creator=dockonaut")
	for _, label := range network.Labels {
		args = append(args, "--label", label)
	}
	for _, option := range network.Options {
		args = append(args, "--opt", option)
	}
	if network.Attachable {
		args = append(args, "--attachable")
	}
	if network.IPV6 {
		args = append(args, "--ipv6")
	}
	if network.Internal {
		args = append(args, "--internal")
	}
	if network.SubNet != "" {
		args = append(args, "--subnet="+network.SubNet)
	}
	if network.IPRange != "" {
		args = append(args, "--ip-range="+network.IPRange)
	}
	if network.Gateway != "" {
		args = append(args, "--gateway="+network.Gateway)
	}

	args = append(args, network.Name)
	_, err := DockerNetworkExec(args...)
	return err
}

func (network *Network) Drop() error {
	_, err := DockerNetworkExec("rm", network.Name)
	return err
}
