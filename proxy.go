package wgcf

import (
	"os"
	"syscall"

	"github.com/KazumiLine/wgcf/wireproxy"
	"golang.zx2c4.com/wireguard/device"
)

func GenerateTunnel(silent bool) (*wireproxy.VirtualTun, error) {
	prof, err := Generate()
	if err != nil {
		return nil, err
	}
	conf, err := wireproxy.ParseConfigRaw(prof)
	if err != nil {
		return nil, err
	}

	os.Stdout = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
	logLevel := device.LogLevelVerbose
	if silent {
		logLevel = device.LogLevelSilent
	}
	return wireproxy.StartWireguard(conf.Device, logLevel)
}
