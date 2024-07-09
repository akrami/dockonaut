package docker

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os/exec"
)

func DockerComposeExec(args ...string) (string, error) {
	args = append([]string{"docker", "compose"}, args...)
	return CommandExecute(args...)
}

func DockerComposeExecScanner(args ...string) (*bufio.Scanner, exec.Cmd) {
	args = append([]string{"docker", "compose"}, args...)
	command := &Command{Args: args}
	return command.Run()
}

func DockerNetworkExec(args ...string) (string, error) {
	args = append([]string{"docker", "network"}, args...)
	return CommandExecute(args...)
}

func DockerVolumeExec(args ...string) (string, error) {
	args = append([]string{"docker", "volume"}, args...)
	return CommandExecute(args...)
}

func GitExec(args ...string) (*bufio.Scanner, exec.Cmd) {
	args = append([]string{"git"}, args...)
	command := &Command{Args: args}
	return command.Run()
}

func CommandExecute(args ...string) (string, error) {
	var cmd exec.Cmd
	if len(args) < 2 {
		cmd = *exec.Command(args[0])
	} else {
		cmd = *exec.Command(args[0], args[1:]...)
	}
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output[:]), errors.Join(errors.New("runtime error"), err)
	}
	return string(output[:]), nil
}

func (command *Command) Run() (*bufio.Scanner, exec.Cmd) {
	var cmd exec.Cmd
	if len(command.Args) < 2 {
		cmd = *exec.Command(command.Args[0])
	} else {
		cmd = *exec.Command(command.Args[0], command.Args[1:]...)
	}

	if command.WorkingDirectory != "" {
		cmd.Dir = command.WorkingDirectory
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	mergedPipe := io.MultiReader(stdout, stderr)
	cmd.Start()

	scanner := bufio.NewScanner(mergedPipe)
	scanner.Split(bufio.ScanLines)

	return scanner, cmd
}

func LogOutput(scanner bufio.Scanner, cmd exec.Cmd) {
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
	cmd.Wait()
}
