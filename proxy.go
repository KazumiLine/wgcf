package wgcf

import (
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

	logLevel := device.LogLevelVerbose
	if silent {
		// os.Stdout = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
		logLevel = device.LogLevelSilent
	}
	return wireproxy.StartWireguard(conf.Device, logLevel)
}
