package aaaafilter

import (
	"sync"

	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/vanguard/acl"
	"github.com/ben-han-cn/vanguard/config"
	"github.com/ben-han-cn/vanguard/core"
	"github.com/ben-han-cn/vanguard/httpcmd"
)

type aaaaFilter struct {
	viewAndAcls map[string][]string
	lock        sync.RWMutex
}

func NewAAAAFilter() *aaaaFilter {
	filter := &aaaaFilter{}
	httpcmd.RegisterHandler(filter, []httpcmd.Command{&PutAAAAFilter{}})
	return filter
}

func (f *aaaaFilter) ReloadConfig(conf *config.VanguardConf) {
	viewAndAcls := make(map[string][]string)
	for _, aaaaConf := range conf.AAAAFilter {
		viewAndAcls[aaaaConf.View] = aaaaConf.Acls
	}
	f.viewAndAcls = viewAndAcls
}

func (f *aaaaFilter) TransferResponse(cli *core.Client) {
	if f.isTransferEnable(cli) == false {
		return
	}

	if cli.Response == nil || cli.Response.Header.Rcode != g53.R_NOERROR {
		return
	}

	f.removeAaaaRecords(cli.Response)
}

func (f *aaaaFilter) removeAaaaRecords(response *g53.Message) {
	removeAaaaInSection(response, g53.AnswerSection)
	removeAaaaInSection(response, g53.AdditionalSection)
}

func removeAaaaInSection(response *g53.Message, section g53.SectionType) {
	rrsets := response.Sections[section]
	var newRRsets []*g53.RRset
	for _, rrset := range rrsets {
		if rrset.Type != g53.RR_AAAA {
			newRRsets = append(newRRsets, rrset)
		}
	}
	response.Sections[section] = newRRsets
}

func (f *aaaaFilter) isTransferEnable(cli *core.Client) bool {
	f.lock.RLock()
	acls, ok := f.viewAndAcls[cli.View]
	if ok == false {
		f.lock.RUnlock()
		return false
	}
	f.lock.RUnlock()

	for _, aclName := range acls {
		if acl.GetAclManager().Find(aclName, cli.IP()) {
			return true
		}
	}

	return false
}
