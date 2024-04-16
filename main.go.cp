package main

import (
	"akrami/dockonaut/src/config"
	"akrami/dockonaut/src/docker"
	"fmt"
)

func main() {
	repos := config.Load("config.json")
	for _, dep := range repos.Dependencies {

		// var extraField map[string]any
		// if dep.Extra != nil {
		// 	extraField = dep.Extra.(map[string]any)
		// }

		if dep.Type == config.DepTypeVolume {
			volume, err := docker.GetVolume(dep.Name)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(volume)
			// jsonOutput, jsonErr := docker.GetVolume(dep.Name)

			// if jsonErr != nil {

			// 	var driver string
			// 	if extraField["labels"] != nil {
			// 		driver = extraField["driver"].(string)
			// 	} else {
			// 		driver = "local"
			// 	}

			// 	var labels []string
			// 	if extraField["labels"] != nil {
			// 		labels = make([]string, len(extraField["labels"].([]interface{})))
			// 		for i, v := range extraField["labels"].([]interface{}) {
			// 			labels[i] = fmt.Sprint(v)
			// 		}
			// 	}

			// 	var options []string
			// 	if extraField["options"] != nil {
			// 		options = make([]string, len(extraField["options"].([]interface{})))
			// 		for i, v := range extraField["options"].([]interface{}) {
			// 			options[i] = fmt.Sprint(v)
			// 		}
			// 	}

			// 	config := docker.Volume{
			// 		Name:    dep.Name,
			// 		Driver:  driver,
			// 		Labels:  labels,
			// 		Options: options,
			// 	}

			// 	_, err := docker.CreateVolume(config)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// } else {
			// 	// the volume is there
			// 	// check if is compatible
			// 	fmt.Println(jsonOutput)
			// }
		} else if dep.Type == config.DepTypeNetwork {
			fmt.Println("network dependency")
		} else {
			panic("dependency type not defined")
		}
	}
}
