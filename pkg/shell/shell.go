package shell

import "os/exec"

func Execute(cmd string) error {
	command := exec.Command("bash", "-c", cmd)
	return command.Run()
}
