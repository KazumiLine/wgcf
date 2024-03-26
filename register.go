package wgcf

import (
	"fmt"
	"log"

	"github.com/KazumiLine/wgcf/cloudflare"
	"github.com/KazumiLine/wgcf/config"
	"github.com/KazumiLine/wgcf/wireguard"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func RegisterAccount(deviceName string, deviceModel string, existingKey string, acceptedTOS bool) error {
	if IsConfigValidAccount() {
		return errors.New("existing account detected")
	}
	if accepted, err := checkTOS(acceptedTOS); err != nil || !accepted {
		return err
	}

	var privateKey *wireguard.Key
	var err error

	if existingKey != "" {
		privateKey, err = wireguard.NewKey(existingKey)
	} else {
		privateKey, err = wireguard.NewPrivateKey()
	}
	if err != nil {
		return err
	}

	device, err := cloudflare.Register(privateKey.Public(), deviceModel)
	if err != nil {
		return err
	}

	viper.Set(config.PrivateKey, privateKey.String())
	viper.Set(config.DeviceId, device.Id)
	viper.Set(config.AccessToken, device.Token)
	viper.Set(config.LicenseKey, device.Account.License)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	ctx := CreateContext()
	_, err = SetDeviceName(ctx, deviceName)
	if err != nil {
		return err
	}
	thisDevice, err := cloudflare.GetSourceDevice(ctx)
	if err != nil {
		return err
	}

	boundDevice, err := cloudflare.UpdateSourceBoundDeviceActive(ctx, true)
	if err != nil {
		return err
	}
	if !boundDevice.Active {
		return errors.New("failed to activate device")
	}

	printDeviceData(thisDevice, boundDevice)
	log.Println("Successfully created Cloudflare Warp account")
	return nil
}

func checkTOS(acceptedTOS bool) (bool, error) {
	if !acceptedTOS {
		fmt.Println("This project is in no way affiliated with Cloudflare")
		fmt.Println("Cloudflare's Terms of Service: https://www.cloudflare.com/application/terms/")
		prompt := promptui.Select{
			Label: "Do you agree?",
			Items: []string{"Yes", "No"},
		}
		if _, result, err := prompt.Run(); err != nil || result != "Yes" {
			return false, err
		}
	}
	return true, nil
}
