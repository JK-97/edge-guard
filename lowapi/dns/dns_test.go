package dns

import "testing"

func TestPasrseIPInTxt(t *testing.T) {
    url := "port30111.version-control.ffffffffffffffffffffffff.master.iotedge"
    ParseIpInTxt(url)
}
