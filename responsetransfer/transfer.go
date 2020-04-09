package responsetransfer

import (
	"github.com/ben-han-cn/vanguard/config"
	"github.com/ben-han-cn/vanguard/core"
	"github.com/ben-han-cn/vanguard/responsetransfer/aaaafilter"
	"github.com/ben-han-cn/vanguard/responsetransfer/hijack"
	"github.com/ben-han-cn/vanguard/responsetransfer/sortlist"
)

const (
	AAAAFilter string = "aaaa_filter"
	Hijack     string = "hijack"
	Sortlist   string = "sortlist"
)

type Transfer interface {
	TransferResponse(*core.Client)
}

type transferAdaptor struct {
	core.DefaultHandler
	transfer Transfer
}

func NewAAAAFilter(conf *config.VanguardConf) core.DNSQueryHandler {
	return newAdaptor(aaaafilter.NewAAAAFilter(), conf)
}

func NewHijack(conf *config.VanguardConf) core.DNSQueryHandler {
	return newAdaptor(hijack.NewHijack(), conf)
}

func NewSortList(conf *config.VanguardConf) core.DNSQueryHandler {
	return newAdaptor(sortlist.NewSortList(), conf)
}

func newAdaptor(t Transfer, conf *config.VanguardConf) core.DNSQueryHandler {
	a := &transferAdaptor{
		transfer: t,
	}
	a.ReloadConfig(conf)
	return a
}

func (a *transferAdaptor) ReloadConfig(conf *config.VanguardConf) {
	config.ReloadConfig(a.transfer, conf)
}

func (a *transferAdaptor) HandleQuery(ctx *core.Context) {
	core.PassToNext(a, ctx)

	a.transfer.TransferResponse(&ctx.Client)
}
