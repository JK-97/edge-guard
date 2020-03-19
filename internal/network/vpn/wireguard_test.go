package vpn

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var testData string = `[Interface]
PrivateKey = gHUQ270a9PAyjGdzjOmWhZmDYLCPPuBRF3XTJjmcRlQ=
Address = 10.209.10.17/32

[Peer]
PublicKey = k7XM3E3H2VXel2KfrBhTUSP11TWXm9VKJ7AQDmLeWlY=
Endpoint = 152.136.212.176:52000
AllowedIPs = 10.209.10.1/24, 172.21.0.0/16
PersistentKeepalive = 21
`

func TestParseWireGuardConfig(t *testing.T) {
	err := os.MkdirAll(path.Join("/etc", "wireguard"), 0755)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(wireguardConfigFile, []byte(testData), 0755)
	if err != nil {
		t.Error(err)
	}
	ip, err := ParseWireGuardConfig()
	if err != nil {
		t.Log(err)
	}
	if ip != "152.136.212.176" {
		t.Error("failed")
	} else {
		t.Log(ip)
	}
	err = os.Remove(wireguardConfigFile)
	if err != nil {
		t.Log(err)
	}
}
