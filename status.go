package wgcf

import (
	"github.com/KazumiLine/wgcf/cloudflare"
	"github.com/pkg/errors"
)

func Status() (*cloudflare.Device, *cloudflare.BoundDevice, error) {
	if !IsConfigValidAccount() {
		return nil, nil, errors.New("no valid account detected")
	}
	ctx := CreateContext()
	thisDevice, err := cloudflare.GetSourceDevice(ctx)
	if err != nil {
		return nil, nil, err
	}
	boundDevice, err := cloudflare.GetSourceBoundDevice(ctx)
	if err != nil {
		return nil, nil, err
	}
	printDeviceData(thisDevice, boundDevice)
	return thisDevice, boundDevice, nil
}
