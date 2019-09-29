package journal_test

import (
	"bytes"
	"jxcore/journal"
	"testing"
)

func TestDeserialize(t *testing.T) {
	blob := `{
	"create_time": 1564329600,
	"files": [
		"kernel.log",
		"supervisord.log",
		"dockerd.log"
	],
	"mode": "kernel",
	"modify_time": 1564394400,
	"rotates": {}
}`

	var meta journal.LogArchiveMeta

	err := meta.Deserialize(bytes.NewBufferString(blob))
	if err != nil {
		t.Error(err)
	}

	if meta.CreateTime != 1564329600 {
		t.Error("Deserialize failed")
	}

	if meta.ModifyTime != 1564394400 {
		t.Error("Deserialize failed")
	}

	if len(meta.Files) != 3 {
		t.Error("Deserialize failed")
	}

	if meta.Mode != "kernel" {
		t.Error("Deserialize failed")
	}
}
