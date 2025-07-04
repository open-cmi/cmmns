package time

func GetTimeZoneList() ([]string, error) {
	return []string{"Asia/Shanghai"}, nil
}

func ApplyTimeZone(tz string) error {
	return nil
}
