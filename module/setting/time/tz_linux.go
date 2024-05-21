package time

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/service/initial"
)

func GetTimeZoneList() ([]string, error) {
	cmd := exec.Command("bash", "-c", "timedatectl list-timezones")
	out, err := cmd.Output()
	if err != nil {
		logger.Errorf("list timezones failed: %s\n", err.Error())
		return []string{"Asia/Shanghai"}, nil
	}
	arrs := strings.Split(string(out), "\n")
	return arrs, nil
}

func SetTimeZone(tz string) error {
	cmdStr := fmt.Sprintf("timedatectl set-timezone %s", tz)
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		logger.Errorf("exec %s command failed: %s\n", cmdStr, err.Error())
		return errors.New("set time zone failed")
	}
	return nil
}

func Init() error {
	m := Get()
	if m == nil {
		m = New()
		m.AutoAdjust = true
		m.NtpServer = "ntp.ubuntu.com"
		m.TimeZone = "Asia/Shanghai"
		return m.Save()
	}
	return nil
}

func init() {
	initial.Register("time-setting", initial.DefaultPriority, Init)
}
