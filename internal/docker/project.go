package docker

import (
	"akrami/dockonaut/internal/engine"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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

func (project *Project) Start(subtextChannel chan<- engine.Subtext) error {
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
			subtextChannel <- engine.Subtext(fmt.Sprintf("copying directory %s to %s", srcDir, getAbsoluteProjectPath(*project)))
			err := copyutil.CopyDir(filesys.MakeFsOnDisk(), srcDir, getAbsoluteProjectPath(*project))
			if err != nil {
				return errors.Join(fmt.Errorf("cannot copy directory %s to %s", srcDir, getAbsoluteProjectPath(*project)), err)
			}
		} else {
			subtextChannel <- engine.Subtext(fmt.Sprintf("cloning %s", project.Name))
			gitArgs := []string{"clone", project.Repository, getAbsoluteProjectPath(*project)}
			if project.Branch != "" {
				gitArgs = append(gitArgs, "-b", project.Branch)
			}
			scanner, cmd := GitExec(gitArgs...)
			LogOutput(*scanner, cmd, subtextChannel)
		}

		for _, script := range project.PreAction {
			subtextChannel <- engine.Subtext(fmt.Sprintf("script '%s'", script))
			command := &Command{
				WorkingDirectory: getAbsoluteProjectPath(*project),
				Args:             strings.Split(script, " "),
			}
			scanner, cmd := command.Run()
			LogOutput(*scanner, cmd, subtextChannel)
		}
	}
	subtextChannel <- engine.Subtext(fmt.Sprintf("starting %s", project.Name))
	scanner, cmd := DockerComposeExecScanner("-f", getAbsoluteComposePath(*project), "up", "-d", "--build", "--pull", "always")
	LogOutput(*scanner, cmd, subtextChannel)
	if firstRun {
		timeout := 30
		for !project.IsRunning() {
			if timeout < 0 {
				return errors.New("timeout waiting for containers to be up")
			}
			time.Sleep(time.Second)
			timeout--
		}
		for _, script := range project.PostAction {
			subtextChannel <- engine.Subtext(fmt.Sprintf("script '%s'", script))
			command := &Command{
				WorkingDirectory: getAbsoluteProjectPath(*project),
				Args:             strings.Split(script, " "),
			}
			scanner, cmd := command.Run()
			LogOutput(*scanner, cmd, subtextChannel)
		}
	}
	return nil
}

func (project *Project) Restart(deep bool, subtextCahnnel chan<- engine.Subtext) error {
	err := project.Stop(!deep, subtextCahnnel)
	if err != nil {
		return errors.Join(errors.New("project stop error"), err)
	}
	if deep {
		err = project.Purge(subtextCahnnel)
		if err != nil {
			return errors.Join(errors.New("project purge error"), err)
		}
	}
	err = project.Start(subtextCahnnel)
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

func (project *Project) Stop(soft bool, subtextCahnnel chan<- engine.Subtext) error {
	command := "down"
	if soft {
		command = "stop"
	}
	_, err := os.Stat(getAbsoluteComposePath(*project))
	if err != nil {
		return nil
	}
	subtextCahnnel <- engine.Subtext(fmt.Sprintf("stopping %s", project.Name))
	_, dockerErr := DockerComposeExec("-f", getAbsoluteComposePath(*project), command)
	if dockerErr != nil {
		subtextCahnnel <- engine.Subtext(fmt.Sprintf("error in stopping %s", project.Name))
		return errors.Join(fmt.Errorf("docker compose %s error", command), dockerErr)
	}
	return nil
}

func (project *Project) Purge(subtextCahnnel chan<- engine.Subtext) error {
	project.Stop(false, subtextCahnnel)
	subtextCahnnel <- engine.Subtext("cleaning workspace")
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
