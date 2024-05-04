package docker

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
)

func DockerComposeExec(args ...string) (string, error) {
	args = append([]string{"docker", "compose"}, args...)
	return CommandExecute(args...)
}

func DockerComposeExecLive(args ...string) (*bufio.Scanner, exec.Cmd) {
	args = append([]string{"docker", "compose"}, args...)
	return CommandExecuteLive(args...)
}

func DockerContainerExec(args ...string) (string, error) {
	args = append([]string{"docker", "container"}, args...)
	return CommandExecute(args...)
}

func DockerNetworkExec(args ...string) (string, error) {
	args = append([]string{"docker", "network"}, args...)
	return CommandExecute(args...)
}

func DockerVolumeExec(args ...string) (string, error) {
	args = append([]string{"docker", "volume"}, args...)
	return CommandExecute(args...)
}

func GitExec(args ...string) (string, error) {
	args = append([]string{"git"}, args...)
	return CommandExecute(args...)
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
	output, err := cmd.Output()
	if err != nil {
		return string(output[:]), errors.Join(errors.New("runtime error"), err)
	}
	return string(output[:]), nil
}

func CommandExecuteLive(args ...string) (*bufio.Scanner, exec.Cmd) {
	var cmd exec.Cmd
	if len(args) < 2 {
		cmd = *exec.Command(args[0])
	} else {
		cmd = *exec.Command(args[0], args[1:]...)
	}
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	mergedPipe := io.MultiReader(stdout, stderr)
	cmd.Start()

	scanner := bufio.NewScanner(mergedPipe)
	scanner.Split(bufio.ScanLines)
	return scanner, cmd
}
