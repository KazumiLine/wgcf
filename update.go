package wgcf

import (
	"log"

	"github.com/KazumiLine/wgcf/cloudflare"
	"github.com/KazumiLine/wgcf/config"
	"github.com/pkg/errors"
)

func UpdateAccount(licenseKey string, deviceName string) error {
	if !IsConfigValidAccount() {
		return errors.New("no account detected")
	}
	ctx := CreateContext()
	ctx.LicenseKey = licenseKey
	thisDevice, err := cloudflare.GetSourceDevice(ctx)
	if err != nil {
		return err
	}
	_, thisDevice, err = ensureLicenseKeyUpToDate(ctx, thisDevice)
	if err != nil {
		return err
	}
	boundDevice, err := cloudflare.GetSourceBoundDevice(ctx)
	if err != nil {
		return err
	}
	if boundDevice.Name == nil || (deviceName != "" && deviceName != *boundDevice.Name) {
		log.Println("Setting device name")
		if _, err := SetDeviceName(ctx, deviceName); err != nil {
			return err
		}
	}
	boundDevice, err = cloudflare.UpdateSourceBoundDeviceActive(ctx, true)
	if err != nil {
		return err
	}
	if !boundDevice.Active {
		return errors.New("failed activating device")
	}
	printDeviceData(thisDevice, boundDevice)
	log.Println("Successfully updated Cloudflare Warp account")
	return nil
}

func ensureLicenseKeyUpToDate(ctx *config.Context, thisDevice *cloudflare.Device) (*cloudflare.Account, *cloudflare.Device, error) {
	if thisDevice.Account.License != ctx.LicenseKey {
		log.Println("Updated license key detected, re-binding device to new account")
		return updateLicenseKey(ctx)
	}
	return nil, thisDevice, nil
}

func updateLicenseKey(ctx *config.Context) (*cloudflare.Account, *cloudflare.Device, error) {
	if _, err := cloudflare.UpdateLicenseKey(ctx); err != nil {
		return nil, nil, err
	}
	account, err := cloudflare.GetAccount(ctx)
	if err != nil {
		return nil, nil, err
	}
	thisDevice, err := cloudflare.GetSourceDevice(ctx)
	if err != nil {
		return nil, nil, err
	}
	if account.License != ctx.LicenseKey {
		return nil, nil, errors.New("failed to update license key")
	}
	return account, thisDevice, nil
}
