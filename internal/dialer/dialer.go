package dialer

import (
	"net"
	"time"
)

func New(dialerTimeout, dialerKeepAlive time.Duration) *net.Dialer {
	dialer := &net.Dialer{
		Timeout:   dialerTimeout,
		KeepAlive: dialerKeepAlive,
	}

	return dialer
}
