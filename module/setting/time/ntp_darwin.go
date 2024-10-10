package time

func SetTimeSetting(req *SettingRequest) error {
	s := Get()
	if s == nil {
		s = New()
	}
	s.NtpServer = req.NtpServer
	s.AutoAdjust = req.AutoAdjust

	err := s.Save()
	return err
}
