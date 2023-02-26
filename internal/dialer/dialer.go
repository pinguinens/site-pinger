package dialer

import (
	"net"
	"time"
)

func New(dialerTimeout, dialerKeepAlive time.Duration) *net.Dialer {
	dialer := &net.Dialer{
		Timeout:   dialerTimeout * time.Second,
		KeepAlive: dialerKeepAlive * time.Second,
	}

	return dialer
}
