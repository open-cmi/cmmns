package device

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DeviceID global var
var DeviceID string = ""

// GetLinuxProductID func
func GetLinuxProductID() string {
	filePth := "/sys/class/dmi/id/product_uuid"
	file, err := os.Open(filePth)
	if err != nil {
		return ""
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	deviceid := strings.Trim(string(data), " \r\n\t")
	return deviceid
}

// GetDarwinProductID func
func GetDarwinProductID() string {
	grep := exec.Command("grep", "IOPlatformSerialNumber")
	ioreg := exec.Command("ioreg", "-l")

	// Get ps's stdout and attach it to grep's stdin.
	pipe, err := ioreg.StdoutPipe()
	if err != nil {
		return ""
	}

	defer pipe.Close()
	grep.Stdin = pipe
	// Run ps first.
	ioreg.Start()

	// Run and get the output of grep.
	output, err := grep.Output()
	if err != nil {
		return ""
	}

	arr := strings.Split(string(output), "=")
	if len(arr) != 2 {
		return ""
	}
	deviceid := strings.Trim(arr[1], " ")
	deviceid = strings.Trim(deviceid, "\n")
	deviceid = strings.Trim(deviceid, "\"")
	return deviceid
}

// GetDeviceID func
func GetDeviceID() (deviceid string) {
	var sys string

	if DeviceID != "" {
		return DeviceID
	}

	sys = runtime.GOOS
	if sys == "darwin" {
		deviceid = GetDarwinProductID()
		// system_profiler SPHardwareDataType | awk '/Hardware UUID:/ {print $NF}'
	} else if sys == "linux" {
		deviceid = GetLinuxProductID()
	} else {
		deviceid = "not support"
	}

	DeviceID = deviceid
	return deviceid
}
