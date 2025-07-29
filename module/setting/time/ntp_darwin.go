package time

import (
	"errors"

	"github.com/open-cmi/gobase/essential/i18n"
)

func SetTimeSetting(req *SettingRequest) error {
	if req.NtpServer == "" {
		return errors.New(i18n.Sprintf("ntp server should not be empty"))
	}

	s := Get()
	if s == nil {
		s = New()
	}
	s.NtpServer = req.NtpServer
	s.AutoAdjust = req.AutoAdjust

	err := s.Save()
	return err
}
