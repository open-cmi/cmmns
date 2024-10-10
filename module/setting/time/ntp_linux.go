package time

import (
	"fmt"
	"os"
	"os/exec"
)

func ChangeNTPServer(server string) error {
	if server != "" {
		wf, err := os.OpenFile("/etc/systemd/timesyncd.conf", os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		wf.WriteString("[Time]\n")
		line := fmt.Sprintf("NTP=%s\n", server)
		wf.WriteString(line)
	}
	cmd := exec.Command("bash", "-c", "systemctl restart systemd-timesyncd")
	err := cmd.Run()
	return err
}

func SetTimeSetting(req *SettingRequest) error {
	s := Get()
	if s == nil {
		s = New()
	}
	s.NtpServer = req.NtpServer
	s.AutoAdjust = req.AutoAdjust
	s.TimeZone = req.TimeZone

	err := ChangeNTPServer(req.NtpServer)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if s.AutoAdjust {
		cmd = exec.Command("bash", "-c", "timedatectl set-ntp true")
	} else {
		cmd = exec.Command("bash", "-c", "timedatectl set-ntp false")
	}
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = SetTimeZone(s.TimeZone)
	if err != nil {
		return err
	}

	err = s.Save()
	return err
}
