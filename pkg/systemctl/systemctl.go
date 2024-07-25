package systemctl

import "os/exec"

func StartService(service string) error {
	cmd := exec.Command("systemctl", "enable", service)
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("systemctl", "start", service)
	return cmd.Run()
}

func StopService(service string) error {
	cmd := exec.Command("systemctl", "stop", service)
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("systemctl", "disable", service)
	return cmd.Run()
}

func StatusService(service string) bool {
	cmd := exec.Command("systemctl", "is-active", service)
	err := cmd.Run()
	return err == nil
}

func RestartService(service string) error {
	cmd := exec.Command("systemctl", "restart", service)
	err := cmd.Run()
	return err
}
