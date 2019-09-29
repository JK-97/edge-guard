package docker

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	since := int64(1568271599)

	st := time.Unix(since, 0).UTC()

	println(st.Format(time.RFC3339))
}
