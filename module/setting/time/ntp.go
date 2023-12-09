package time

func SetNtpSetting(req *SettingRequest) error {
	s := Get()
	if s == nil {
		s = New()
	}
	s.NtpServer = req.NtpServer
	s.AutoAdjust = req.AutoAdjust
	err := s.Save()
	return err
}

func GetNtpSetting() *Setting {
	s := Get()
	if s == nil {
		s = New()
	}
	return s
}

func Adjust() error {
	return nil
}
