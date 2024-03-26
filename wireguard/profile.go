package wireguard

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

var profileTemplate = `[Interface]
PrivateKey = {{ .PrivateKey }}
Address = {{ .Address1 }}/32
Address = {{ .Address2 }}/128
DNS = 1.1.1.1, 1.0.0.1, 2606:4700:4700::1111, 2606:4700:4700::1001
MTU = 1280
[Peer]
PublicKey = {{ .PublicKey }}
AllowedIPs = 0.0.0.0/0
AllowedIPs = ::/0
Endpoint = {{ .Endpoint }}
`

type Profile struct {
	PrivateKey string
	Address1   string
	Address2   string
	PublicKey  string
	Endpoint   string
}

func GenerateProfile(data *Profile) (string, error) {
	t, err := template.New("").Parse(profileTemplate)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	if err := t.Execute(&result, data); err != nil {
		return "", err
	}
	return result.String(), nil
}

func (p *Profile) Save(profileFile string) error {
	profileData, err := GenerateProfile(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(profileFile, []byte(profileData), 0600)
}
