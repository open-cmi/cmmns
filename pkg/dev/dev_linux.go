package dev

import (
	"io/ioutil"
	"os"
	"strings"
)

var deviceIDFiles []string = []string{
	"/sys/class/dmi/id/product_uuid",
	"/sys/block/mmcblk0/device/serial",
}

// GetDeviceID func
func GetDeviceID() string {
	for _, filep := range deviceIDFiles {
		file, err := os.Open(filep)
		if err != nil {
			continue
		}

		data, _ := ioutil.ReadAll(file)
		deviceid := strings.Trim(string(data), " \r\n\t ")
		file.Close()
		return deviceid
	}
	return ""
}
