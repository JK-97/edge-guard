package rfile

import "testing"

func TestGetSuffix(t *testing.T) {
	tt := int64(1568865226)
	excepted := "2019-09-19.log"
	actuly := GetSuffix(tt)

	if actuly != excepted {
		t.Error(excepted, actuly)
	}
}
