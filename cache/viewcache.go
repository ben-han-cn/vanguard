package cache

import (
	"github.com/ben-han-cn/g53"
)

type ViewCache struct {
	positiveCache *MessageCache
	negativeCache *MessageCache
}

func newViewCache(cap int) *ViewCache {
	return &ViewCache{
		positiveCache: newMessageCache(cap),
		negativeCache: newMessageCache(cap),
	}
}

func (c *ViewCache) ResetCapacity(cap int) {
	c.positiveCache.ResetCapacity(cap)
	c.negativeCache.ResetCapacity(cap)
}

func (c *ViewCache) Get(req *g53.Message) (*g53.Message, bool) {
	if msg, ok := c.positiveCache.Get(req); ok {
		return msg, true
	} else if msg, ok := c.negativeCache.Get(req); ok {
		return msg, true
	} else {
		return nil, false
	}
}

func (c *ViewCache) Add(msg *g53.Message) {
	if msg.Header.ANCount > 0 {
		c.positiveCache.Add(msg)
	} else {
		c.negativeCache.Add(msg)
	}
}

func (c *ViewCache) Remove(name *g53.Name, typ g53.RRType) bool {
	if found := c.positiveCache.Remove(name, typ); found {
		return true
	} else {
		return c.negativeCache.Remove(name, typ)
	}
}

func (c *ViewCache) GetDeepestNS(name *g53.Name) (*g53.RRset, bool) {
	return c.positiveCache.GetDeepestNS(name)
}

func (c *ViewCache) GetRRset(name *g53.Name, typ g53.RRType) (*g53.RRset, bool) {
	return c.positiveCache.GetRRset(name, typ)
}
