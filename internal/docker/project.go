package docker

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/copyutil"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func (project *Project) IsRunning() bool {
	_, err := os.Stat(getAbsoluteProjectPath(*project))
	if err != nil {
		return false
	}
	result, err := DockerComposeExec("-f", getAbsoluteComposePath(*project), "ps", "--all", "--format", "json")
	if err != nil || len(result) == 0 {
		return false
	}
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		var container Container
		json.Unmarshal([]byte(scanner.Text()), &container)
		if container.State == "running" {
			// if at least one of containers is running we consider project as running
			return true
		}
	}

	return false
}

func (project *Project) Start() error {
	firstRun := false
	_, err := os.Stat(getAbsoluteProjectPath(*project))
	if err != nil {
		firstRun = true
		os.MkdirAll(getAbsoluteProjectPath(*project), os.ModePerm)
		if project.Local != "" {
			var srcDir string
			if strings.HasPrefix(project.Local, "/") {
				srcDir = project.Local
			} else {
				root, _ := os.Getwd()
				srcDir = root + "/" + project.Local
			}
			copyutil.CopyDir(filesys.MakeFsOnDisk(), srcDir, getAbsoluteProjectPath(*project))
		} else {
			_, execErr := GitExec("clone", project.Repository, getAbsoluteProjectPath(*project))
			if execErr != nil {
				return errors.Join(errors.New("git exec error"), execErr)
			}
		}

		for _, script := range project.PreAction {
			log.Println("running:", script)
			output, err := CommandExecute(strings.Split(script, " ")...)
			if err != nil {
				return errors.Join(errors.New("cannot run script: "+script), err)
			}
			log.Println("output:", output)
		}
	}
	scanner, cmd := DockerComposeExecLive("-f", getAbsoluteComposePath(*project), "up", "-d", "--build", "--pull", "always")
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
	cmd.Wait()
	if firstRun {
		for _, script := range project.PostAction {
			log.Println("running:", script)
			output, err := CommandExecute(strings.Split(script, " ")...)
			if err != nil {
				return errors.Join(errors.New("cannot run script: "+script), err)
			}
			log.Println("output:", output)
		}
	}
	return nil
}

func (project *Project) Restart(deep bool) error {
	err := project.Stop(false)
	if err != nil {
		return errors.Join(errors.New("project stop error"), err)
	}
	if deep {
		err = project.Purge()
		if err != nil {
			return errors.Join(errors.New("project purge error"), err)
		}
	}
	err = project.Start()
	if err != nil {
		return errors.Join(errors.New("project start error"), err)
	}
	return nil
}

func (project *Project) GetContainers() ([]Container, error) {
	containers := make([]Container, 0)
	result, err := DockerComposeExec("-f", getAbsoluteComposePath(*project), "ps", "--all", "--format", "json")
	if err != nil || len(result) == 0 {
		return containers, errors.Join(errors.New("can not list containers"), err)
	}
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		var container Container
		json.Unmarshal([]byte(scanner.Text()), &container)
		containers = append(containers, container)
	}
	return containers, nil
}

func (project *Project) Stop(soft bool) error {
	command := "down"
	if soft {
		command = "stop"
	}
	_, dockerErr := DockerComposeExec("-f", getAbsoluteComposePath(*project), command)
	if dockerErr != nil {
		return errors.Join(errors.New("docker compose "+command+" error"), dockerErr)
	}
	return nil
}

func (project *Project) Purge() error {
	err := project.Stop(false)
	if err != nil {
		return errors.Join(errors.New("can not stop project"), err)
	}
	return os.RemoveAll(getAbsoluteProjectPath(*project))
}

func getAbsoluteComposePath(project Project) string {
	root, _ := os.Getwd()
	return root + "/workspace/" + project.Name + "/" + project.Path
}

func getAbsoluteProjectPath(project Project) string {
	root, _ := os.Getwd()
	return root + "/workspace/" + project.Name
}
