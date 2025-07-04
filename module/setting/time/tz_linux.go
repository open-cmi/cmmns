package time

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
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

func ApplyTimeZone(tz string) error {
	cmdStr := fmt.Sprintf("timedatectl set-timezone %s", tz)
	cmd := exec.Command("bash", "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		logger.Errorf("exec %s command failed: %s\n", cmdStr, err.Error())
		return errors.New("set time zone failed")
	}
	return nil
}
