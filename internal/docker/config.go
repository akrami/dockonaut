package docker

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Load(configFile string) (RepoList, error) {
	jsonFile, err := os.Open(configFile)
	defer jsonFile.Close()
	if err != nil {
		return RepoList{}, errors.Join(errors.New("cannot open config file"), err)
	}
	byteValue, _ := io.ReadAll(jsonFile)
	var result RepoList
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return RepoList{}, errors.Join(errors.New("wrong json format"), err)
	}
	return result, nil
}

func Start(repos RepoList) error {
	for _, volume := range repos.Dependency.Volume {
		log.Println("checking volume", volume.Name)
		errVol := volume.Load()
		if errVol != nil {
			log.Println("creting volume", volume.Name)
			errCreate := volume.Create()
			if errCreate != nil {
				return errors.Join(errors.New("error creating volume "+volume.Name), errCreate)
			}
		}
	}

	for _, network := range repos.Dependency.Network {
		log.Println("checking network", network.Name)
		errNet := network.Load()
		if errNet != nil {
			log.Println("creating network", network.Name)
			errCreate := network.Create()
			if errCreate != nil {
				return errors.Join(errors.New("error creating network "+network.Name), errCreate)
			}
		}
	}

	for _, script := range repos.Dependency.Script {
		log.Println("running:", script)
		args := strings.Split(script, " ")
		output, err := CommandExecute(args...)
		if err != nil {
			return errors.Join(errors.New("executing error: "+script), err)
		}
		log.Println("output:", output)
	}

	for _, project := range repos.Projects {
		log.Println("starting project", project.Name)
		err := startProject(project, repos)
		if err != nil {
			return errors.Join(errors.New("error starting repo with "+project.Name+" project"), err)
		}
	}
	return nil
}

func startProject(project Project, repos RepoList) error {
	if !project.IsRunning() {
		for _, depName := range project.Depends {
			depProject, depError := findProject(repos, depName)
			if depError != nil {
				return depError
			}
			if !depProject.IsRunning() {
				// TODO prevent loop
				err := startProject(depProject, repos)
				if err != nil {
					return errors.Join(errors.New("starting dep project failed: "+depProject.Name), err)
				}
			}
		}
		return project.Start()
	}
	return nil
}

func findProject(repos RepoList, name string) (Project, error) {
	for _, project := range repos.Projects {
		if project.Name == name {
			return project, nil
		}
	}
	return Project{}, errors.New("project " + name + " not found")
}

func Stop(repos RepoList) error {
	for _, project := range repos.Projects {
		err := project.Stop(false)
		if err != nil {
			return errors.Join(errors.New("cannot stop "+project.Name), err)
		}
	}
	return nil
}

func Restart(repos RepoList) error {
	for _, project := range repos.Projects {
		err := project.Restart(true)
		if err != nil {
			return errors.Join(errors.New("cannot restart "+project.Name), err)
		}
	}
	return nil
}

func Purge(repos RepoList) error {
	err := Stop(repos)
	if err != nil {
		return err
	}
	root, _ := os.Getwd()
	contents, err := filepath.Glob(root + "/workspace/*")
	if err != nil {
		return err
	}
	for _, path := range contents {
		if !strings.HasPrefix(path, ".") {
			os.RemoveAll(path)
		}
	}
	return nil
}
