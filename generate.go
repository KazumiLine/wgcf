package wgcf

import (
	"log"

	"github.com/KazumiLine/wgcf/cloudflare"
	"github.com/KazumiLine/wgcf/config"
	"github.com/KazumiLine/wgcf/wireguard"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func generateProfile(profileFile string) error {
	if !IsConfigValidAccount() {
		return errors.New("no account detected")
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

	profile, err := wireguard.NewProfile(&wireguard.ProfileData{
		PrivateKey: viper.GetString(config.PrivateKey),
		Address1:   thisDevice.Config.Interface.Addresses.V4,
		Address2:   thisDevice.Config.Interface.Addresses.V6,
		PublicKey:  thisDevice.Config.Peers[0].PublicKey,
		Endpoint:   thisDevice.Config.Peers[0].Endpoint.Host,
	})
	if err != nil {
		return err
	}
	if err := profile.Save(profileFile); err != nil {
		return err
	}

	PrintDeviceData(thisDevice, boundDevice)
	log.Println("Successfully generated WireGuard profile:", profileFile)
	return nil
}
