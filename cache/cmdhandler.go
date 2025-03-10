package cache

import (
	"fmt"

	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/vanguard/httpcmd"
	"github.com/ben-han-cn/vanguard/resolver/auth/zone"
)

type CleanCache struct {
}

func (c *CleanCache) String() string {
	return "name: clean all cache"
}

type CleanViewCache struct {
	View string `json:"view_name"`
}

func (c *CleanViewCache) String() string {
	return fmt.Sprintf("name: clean cache and params:{view:%s}", c.View)
}

type CleanRRsetsCache struct {
	View string `json:"view_name"`
	Name string `json:"domain_name"`
}

func (c *CleanRRsetsCache) String() string {
	return fmt.Sprintf("name: clear cache and params:{name:%s, view:%s}", c.Name, c.View)
}

type CleanDomainCache struct {
	Name string `json:"domain_name"`
}

func (c *CleanDomainCache) String() string {
	return fmt.Sprintf("name: clear cache and params:{name:%s}", c.Name)
}

type GetDomainCache struct {
	Name string `json:"domain_name"`
	Type string `json:"type"`
}

func (g *GetDomainCache) String() string {
	return fmt.Sprintf("name: get rrsets from cache and params:{name:%s, type:%s}", g.Name, g.Type)
}

type GetMessageCache struct {
	View string `json:"view_name"`
	Name string `json:"domain_name"`
	Type string `json:"type"`
}

func (g *GetMessageCache) String() string {
	return fmt.Sprintf("name: get rrsets from cache and params:{name:%s, type:%s, view:%s}", g.Name, g.Type, g.View)
}

func (cache *Cache) HandleCmd(cmd httpcmd.Command) (interface{}, *httpcmd.Error) {
	switch c := cmd.(type) {
	case *CleanCache:
		return cache.cleanAll()
	case *CleanViewCache:
		return cache.cleanInView(c.View)
	case *CleanRRsetsCache:
		return cache.cleanDomainInView(c.View, c.Name)
	case *CleanDomainCache:
		return cache.cleanDomain(c.Name)
	case *GetDomainCache:
		return cache.getDomainCache(c.Name, c.Type)
	case *GetMessageCache:
		return cache.getMessageCacheInView(c.View, c.Name, c.Type)
	default:
		panic("shouldn't be here")
	}
}

func (c *Cache) cleanAll() (interface{}, *httpcmd.Error) {
	return nil, nil
}

func (c *Cache) cleanInView(view string) (interface{}, *httpcmd.Error) {
	return nil, nil
}

func (c *Cache) cleanDomainInView(view string, name string) (interface{}, *httpcmd.Error) {
	return c.cleanRRsetsCache(view, name, zone.SupportRRTypes)
}

func (c *Cache) cleanDomain(name string) (interface{}, *httpcmd.Error) {
	for view, _ := range c.cache {
		if code, err := c.cleanRRsetsCache(view, name, zone.SupportRRTypes); err != nil {
			return code, err
		}
	}
	return nil, nil
}

func (c *Cache) cleanRRsetsCache(view string, name string, types []g53.RRType) (interface{}, *httpcmd.Error) {
	return nil, nil
}

type RRInCache struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
	Ttl   int    `json:"ttl"`
	Rdata string `json:"rdata"`
}

func (c *Cache) getDomainCache(name, typ string) (interface{}, *httpcmd.Error) {
	qtype, err := g53.TypeFromString(typ)
	if err != nil {
		return nil, httpcmd.ErrUnknownRRType.AddDetail(typ)
	}

	types := []g53.RRType{qtype}
	if qtype == g53.RR_ANY {
		types = zone.SupportRRTypes
	}

	var all []RRInCache
	for view, _ := range c.cache {
		for _, qtype := range types {
			if rrs, err := c.getMessageCacheInView(view, name, qtype.String()); err == nil {
				rrs_ := rrs.([]RRInCache)
				if len(rrs_) > 0 {
					all = append(all, rrs_...)
				}
			}
		}
	}
	return all, nil
}

func (c *Cache) getMessageCacheInView(view, name, typ string) (interface{}, *httpcmd.Error) {
	qtype, err := g53.TypeFromString(typ)
	if err != nil {
		return nil, httpcmd.ErrUnknownRRType.AddDetail(typ)
	}

	qname, err := g53.NameFromString(name)
	if err != nil {
		return nil, httpcmd.ErrInvalidName.AddDetail(err.Error())
	}

	if qtype == g53.RR_ANY {
		var rrsets []RRInCache
		for _, t := range zone.SupportRRTypes {
			if rrs, err := c.getSingleMessageCache(view, qname, t); err == nil {
				rrsets = append(rrsets, rrs...)
			}
		}
		return rrsets, nil
	} else {
		return c.getSingleMessageCache(view, qname, qtype)
	}
}

func (c *Cache) getSingleMessageCache(view string, qname *g53.Name, qtype g53.RRType) ([]RRInCache, *httpcmd.Error) {
	return nil, nil
}
