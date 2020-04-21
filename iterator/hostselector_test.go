package iterator

import (
	"net"
	"testing"
	"time"

	ut "github.com/ben-han-cn/cement/unittest"
	"github.com/ben-han-cn/vanguard/logger"
)

func TestRttBasedHostSelector(t *testing.T) {
	logger.UseDefaultLogger("debug")

	selector := newRttBasedHostSelector(100)
	host1 := Host(net.ParseIP("1.1.1.1"))
	host2 := Host(net.ParseIP("2.2.2.2"))
	selector.SetRtt(host1, 10*time.Second)
	selector.SetRtt(host2, 11*time.Second)
	target, succeed := selector.SelectHost([]Host{host1, host2})
	ut.Assert(t, succeed, "")
	ut.Equal(t, target, host1)

	selector.SetRtt(host1, 14*time.Second)
	target, succeed = selector.SelectHost([]Host{host1, host2})
	ut.Assert(t, succeed, "")
	ut.Equal(t, target, host2)

	for i := 0; i < MaxTimeoutCount+1; i++ {
		selector.SetTimeout(host1, 20*time.Second)
		selector.SetTimeout(host2, 20*time.Second)
	}
	_, succeed = selector.SelectHost([]Host{host1, host2})
	ut.Assert(t, !succeed, "")
}
