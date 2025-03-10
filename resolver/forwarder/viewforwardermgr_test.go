package forwarder

import (
	"testing"

	ut "github.com/ben-han-cn/cement/unittest"
	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/vanguard/config"
	"github.com/ben-han-cn/vanguard/httpcmd"
	"github.com/ben-han-cn/vanguard/logger"
	view "github.com/ben-han-cn/vanguard/viewselector"
)

func TestViewFwderCreation(t *testing.T) {
	logger.UseDefaultLogger("error")

	var conf config.VanguardConf
	view.NewSelectorMgr(&conf)
	viewFwderMgr := NewViewFwderMgr(&conf)
	viewFwderMgr.ReloadConfig(&conf)
	err := viewFwderMgr.addForwardZone([]ForwardZoneParam{
		{"default", "a.cn", []string{"1.1.1.1:5555"}, "Order"},
		{"default", "b.cn", []string{"1.1.1.1:4444"}, "Order"},
		{"default", "c.cn", []string{"1.1.1.1:5555"}, "Order"},
	})
	ut.Equal(t, err, (*httpcmd.Error)(nil))
	fwder1 := viewFwderMgr.GetFwder("default", g53.NameFromStringUnsafe("a.cn"))
	fwder2 := viewFwderMgr.GetFwder("default", g53.NameFromStringUnsafe("b.cn"))
	fwder3 := viewFwderMgr.GetFwder("default", g53.NameFromStringUnsafe("c.cn"))
	ut.Assert(t, fwder1 != fwder2, "")
	ut.Equal(t, fwder1, fwder3)
}
