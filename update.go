package wgcf

import (
	"fmt"
	"log"

	"github.com/KazumiLine/wgcf/cloudflare"
	"github.com/KazumiLine/wgcf/config"
	"github.com/KazumiLine/wgcf/wireguard"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

func GenerateWARPPlus() (licenseKey string, err error) {
	privateKey, err := wireguard.NewPrivateKey()
	device, err := cloudflare.Register(privateKey.Public(), "PC")
	viper.Set(config.PrivateKey, privateKey.String())
	viper.Set(config.DeviceId, device.Id)
	viper.Set(config.AccessToken, device.Token)
	viper.Set(config.LicenseKey, device.Account.License)
	viper.WriteConfig()
	account1 := CreateContext()
	privateKey, err = wireguard.NewPrivateKey()
	device, err = cloudflare.Register(privateKey.Public(), "PC")
	viper.Set(config.PrivateKey, privateKey.String())
	viper.Set(config.DeviceId, device.Id)
	viper.Set(config.AccessToken, device.Token)
	viper.Set(config.LicenseKey, device.Account.License)
	viper.WriteConfig()
	account2 := CreateContext()
	fmt.Println(cloudflare.UpdateReferrer(account1, account2.DeviceId))
	oldLicenseKey := account1.LicenseKey
	account1.LicenseKey = "hW17X52Z-1hE542mf-e185pLr6"
	updateLicenseKey(account1)
	account1.LicenseKey = oldLicenseKey
	updateLicenseKey(account1)
	fmt.Println(cloudflare.GetAccount(account1))
	return oldLicenseKey, nil
}
