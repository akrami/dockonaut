package docker

import (
	"akrami/dockonaut/internal/engine"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Load(configFile string) (Config, error) {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return Config{}, errors.Join(errors.New("cannot open config file"), err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	err = json.Unmarshal([]byte(byteValue), &config)
	if err != nil {
		return Config{}, errors.Join(errors.New("wrong json format"), err)
	}
	return config, nil
}

func (config *Config) sort() error {
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	projectMap := make(map[string]Project)

	for _, project := range config.Projects {
		projectMap[project.Name] = project
		if _, exists := graph[project.Name]; !exists {
			graph[project.Name] = []string{}
		}
		for _, dep := range project.Depends {
			graph[dep] = append(graph[dep], project.Name)
			inDegree[project.Name]++
		}
	}

	queue := []string{}
	for name := range projectMap {
		if inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}

	var sortedProjects []Project
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		sortedProjects = append(sortedProjects, projectMap[current])

		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(sortedProjects) != len(config.Projects) {
		return errors.New("cycle detected in the projects dependency graph")
	}
	config.Projects = sortedProjects

	return nil
}

func (config *Config) Start(headerChannel chan<- engine.Header, subtextChannel chan<- engine.Subtext) error {
	for _, volume := range config.Dependency.Volumes {
		headerChannel <- engine.Header(fmt.Sprintf("Volume %s", volume.Name))
		subtextChannel <- engine.Subtext(fmt.Sprintf("checking if volume %s already exists", volume.Name))
		if err := volume.Load(); err != nil {
			subtextChannel <- engine.Subtext(fmt.Sprintf("creating volume %s", volume.Name))
			if errCreate := volume.Create(); errCreate != nil {
				subtextChannel <- engine.Subtext(fmt.Sprintf("error creating volume %s", volume.Name))
				return errors.Join(errors.New("Error creating volume "+volume.Name), errCreate)
			}
		}
	}

	for _, network := range config.Dependency.Networks {
		headerChannel <- engine.Header(fmt.Sprintf("Network %s", network.Name))
		subtextChannel <- engine.Subtext(fmt.Sprintf("checking if network %s already exists", network.Name))
		if err := network.Load(); err != nil {
			subtextChannel <- engine.Subtext(err.Error())
			subtextChannel <- engine.Subtext(fmt.Sprintf("creating network %s", network.Name))
			if errCreate := network.Create(); errCreate != nil {
				subtextChannel <- engine.Subtext(fmt.Sprintf("error creating network %s", network.Name))
				return errors.Join(errors.New("error creating network "+network.Name), errCreate)
			}
		}
	}

	for _, script := range config.Dependency.Scripts {
		headerChannel <- engine.Header(fmt.Sprintf("Script '%s'", script))
		args := strings.Split(script, " ")
		command := &Command{Args: args}
		scanner, cmd := command.Run()
		LogOutput(*scanner, cmd, subtextChannel)
	}

	headerChannel <- engine.Header("Sorting projects base on dependencies")
	if err := config.sort(); err != nil {
		subtextChannel <- engine.Subtext("error sorting projects")
		return err
	}

	for _, project := range config.Projects {
		headerChannel <- engine.Header(fmt.Sprintf("Project %s", project.Name))
		if !project.IsRunning() {
			if err := project.Start(subtextChannel); err != nil {
				subtextChannel <- engine.Subtext(fmt.Sprintf("error starting %s", project.Name))
				return errors.Join(fmt.Errorf("error starting repo for %s", project.Name), err)
			}
		}
	}
	return nil
}

func (config Config) Stop(headerChannel chan<- engine.Header, subtextChannel chan<- engine.Subtext) error {
	for _, project := range config.Projects {
		headerChannel <- engine.Header(fmt.Sprintf("Project %s", project.Name))
		err := project.Stop(false, subtextChannel)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot stop %s", project.Name), err)
		}
	}
	return nil
}

func (config Config) Restart(headerChannel chan<- engine.Header, subtextChannel chan<- engine.Subtext) error {
	for _, project := range config.Projects {
		headerChannel <- engine.Header(fmt.Sprintf("Restaring %s", project.Name))
		err := project.Restart(false, subtextChannel)
		if err != nil {
			return errors.Join(errors.New("cannot restart "+project.Name), err)
		}
	}
	return nil
}

func (config Config) Purge(headerChannel chan<- engine.Header, subtextChannel chan<- engine.Subtext) error {
	err := config.Stop(headerChannel, subtextChannel)
	if err != nil {
		return err
	}
	root, _ := os.Getwd()
	headerChannel <- engine.Header("Cleaning workspace")
	contents, err := filepath.Glob(root + "/workspace/*")
	if err != nil {
		return err
	}
	for _, path := range contents {
		if !strings.HasPrefix(path, ".") {
			subtextChannel <- engine.Subtext(fmt.Sprintf("deleting %s directory", path))
			os.RemoveAll(path)
		}
	}
	return nil
}
