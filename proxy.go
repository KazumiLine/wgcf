package wgcf

import (
	"github.com/KazumiLine/wgcf/wireproxy"
	"golang.zx2c4.com/wireguard/device"
)

func GenerateTunnel() (*wireproxy.VirtualTun, error) {
	prof, err := Generate()
	if err != nil {
		return nil, err
	}
	conf, err := wireproxy.ParseConfigRaw(prof)
	if err != nil {
		return nil, err
	}
	return wireproxy.StartWireguard(conf.Device, device.LogLevelSilent)
}
