package wgcf

import (
	"github.com/KazumiLine/wgcf/cloudflare"
	"github.com/pkg/errors"
)

func status() error {
	if !IsConfigValidAccount() {
		return errors.New("no valid account detected")
	}

	ctx := CreateContext()
	thisDevice, err := cloudflare.GetSourceDevice(ctx)
	if err != nil {
		return err
	}
	boundDevice, err := cloudflare.GetSourceBoundDevice(ctx)
	if err != nil {
		return err
	}

	PrintDeviceData(thisDevice, boundDevice)
	return nil
}
