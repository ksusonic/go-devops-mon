package agent

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getCurrentIps(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Well, seems like your machine has problems with net: %v", r)
		}
	}()
	got, err := getCurrentIps()
	assert.NoError(t, err)
	assert.NotEmptyf(t, got, "expected machine to have any ip address")
	assert.Equal(t, net.ParseIP("127.0.0.1"), got[0], "expected first interface to be 127.0.0.1")
}
