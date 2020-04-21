package iterator

import (
	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/vanguard/cache"
)

type DelegationPoint struct {
	zone          *g53.Name
	serverAndHost map[string][]Host
	probedServer  []*g53.Name
	lameHosts     []Host
}

func newDelegationPoint(zone *g53.Name, serverAndHost map[string][]Host) *DelegationPoint {
	return &DelegationPoint{
		zone:          zone,
		serverAndHost: serverAndHost,
	}
}

func newDelegationPointFromReferralResponse(zone *g53.Name, referral *g53.Message) *DelegationPoint {
	ns := referral.GetSection(g53.AuthSection)
	if len(ns) != 1 {
		panic("referral response doesn't has one ns rrset")
	}

	dp := newDelegationPointFromNS(ns[0], nil)
	for _, glue := range referral.GetSection(g53.AdditionalSection) {
		dp.addGlue(glue)
	}
	return dp
}

func newDelegationPointFromNS(ns *g53.RRset, glues []*g53.RRset) *DelegationPoint {
	serverAndHost := make(map[string][]Host)
	for _, rdata := range ns.Rdatas {
		serverAndHost[rdata.(*g53.NS).Name.String(false)] = nil
	}

	dp := newDelegationPoint(ns.Name, serverAndHost)
	for _, glue := range glues {
		dp.addGlue(glue)
	}
	return dp
}

func (dp *DelegationPoint) addGlue(glue *g53.RRset) {
	key := glue.Name.String(false)
	if hosts, ok := dp.serverAndHost[key]; ok {
		for _, rdata := range glue.Rdatas {
			if a, ok := rdata.(*g53.A); ok {
				hosts = append(hosts, Host(a.Host))
			} else if aaaa, ok := rdata.(*g53.AAAA); ok {
				hosts = append(hosts, Host(aaaa.Host))
			} else {
				panic("glue isn't a or aaaa")
			}
		}
		dp.serverAndHost[key] = hosts
	}
}

func newDelegationPointFromCache(qname *g53.Name, viewCache *cache.ViewCache) (*DelegationPoint, bool) {
	ns, ok := viewCache.GetDeepestNS(qname)
	if !ok {
		return nil, false
	}

	var glues []*g53.RRset
	for _, rdata := range ns.Rdatas {
		glue, ok := viewCache.GetRRset(rdata.(*g53.NS).Name, g53.RR_A)
		if ok {
			glues = append(glues, glue)
		}
	}
	return newDelegationPointFromNS(ns, glues), true
}
