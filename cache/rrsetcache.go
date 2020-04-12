package cache

import (
	"container/list"
	"time"

	"github.com/ben-han-cn/g53"
)

type TrustLevel int

const (
	AdditionalWithoutAA TrustLevel = 0
	AuthorityWithoutAA  TrustLevel = 1
	AdditionalWithAA    TrustLevel = 2
	AnswerWithoutAA     TrustLevel = 3
	AuthorityWithAA     TrustLevel = 4
	AnswerWithAA        TrustLevel = 5
)

type RRsetEntry struct {
	keyHash      uint64
	conflictHash uint64
	rrset        *g53.RRset
	trustLevel   TrustLevel
	expireTime   time.Time
}

func (e *RRsetEntry) IsExpire() bool {
	return e.expireTime.Before(time.Now())
}

type RRsetCache struct {
	cap  int
	data map[uint64]*list.Element
	ll   *list.List
}

func newRRsetCache(cap int) *RRsetCache {
	return &RRsetCache{
		ll:   list.New(),
		data: make(map[uint64]*list.Element),
		cap:  cap,
	}
}

func (c *RRsetCache) add(es []RRsetEntry) {
	for i, e := range es {
		if elem, ok := c.data[e.keyHash]; ok {
			oe := elem.Value.(*RRsetEntry)
			if !oe.IsExpire() && e.trustLevel < oe.trustLevel {
				return
			}
			c.ll.MoveToFront(elem)
			elem.Value = &es[i]
		} else if c.ll.Len() < c.cap {
			elem := c.ll.PushFront(&es[i])
			c.data[e.keyHash] = elem
		} else {
			//reuse last elem
			elem := c.ll.Back()
			oe := elem.Value.(*RRsetEntry)
			delete(c.data, oe.keyHash)
			*oe = e
			c.data[e.keyHash] = elem
			c.ll.MoveToFront(elem)
		}
	}
}

func (c *RRsetCache) get(keyHash, conflictHash uint64) (*g53.RRset, bool) {
	if elem, hit := c.data[keyHash]; hit {
		e := elem.Value.(*RRsetEntry)
		now := time.Now()
		if e.conflictHash == conflictHash && e.expireTime.After(now) {
			c.ll.MoveToFront(elem)
			rrset := *e.rrset
			rrset.Ttl = g53.RRTTL(e.expireTime.Sub(now).Seconds())
			return &rrset, true
		}
	}
	return nil, false
}

func (c *RRsetCache) remove(keyHash, conflictHash uint64) bool {
	if elem, hit := c.data[keyHash]; hit {
		e := elem.Value.(*RRsetEntry)
		if e.conflictHash == conflictHash {
			delete(c.data, keyHash)
			c.ll.Remove(elem)
			return true
		}
	}
	return false
}
