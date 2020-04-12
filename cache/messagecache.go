package cache

import (
	"container/list"
	"math"
	"sync"
	"time"

	"github.com/ben-han-cn/g53"
)

const (
	RRsetVsMessageRatio = 5
)

type RRsetHash struct {
	keyHash      uint64
	conflictHash uint64
}

type MessageEntry struct {
	keyHash         uint64
	conflictHash    uint64
	rcode           g53.Rcode
	answerCount     uint16
	authorityCount  uint16
	additionalCount uint16
	rrsets          []RRsetHash
	expireTime      time.Time
}

func (e *MessageEntry) IsExpire() bool {
	return e.expireTime.Before(time.Now())
}

type MessageCache struct {
	cap        int
	data       map[uint64]*list.Element
	ll         *list.List
	mu         sync.Mutex
	rrsetCache *RRsetCache
}

func newMessageCache(cap int) *MessageCache {
	return &MessageCache{
		ll:         list.New(),
		data:       make(map[uint64]*list.Element),
		cap:        cap,
		rrsetCache: newRRsetCache(cap * RRsetVsMessageRatio),
	}
}

func (c *MessageCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.data)
}

func (c *MessageCache) Get(req *g53.Message) (*g53.Message, bool) {
	c.mu.Lock()

	keyHash, conflictHash := HashQuery(req.Question.Name, req.Question.Type)
	me, found := c.get(keyHash, conflictHash)
	if !found {
		c.mu.Unlock()
		return nil, false
	}

	var rrsets []*g53.RRset
	rrsetCount := me.answerCount + me.authorityCount + me.additionalCount
	if rrsetCount > 0 {
		rrsets = make([]*g53.RRset, rrsetCount)
		for i, hash := range me.rrsets {
			rrset, found := c.rrsetCache.get(hash.keyHash, hash.conflictHash)
			if !found {
				c.remove(keyHash, conflictHash)
				c.mu.Unlock()
				return nil, false
			}
			rrsets[i] = rrset
		}
	}
	c.mu.Unlock()

	resp := req.MakeResponse()
	j := 0
	for i := uint16(0); i < me.answerCount; i++ {
		resp.AddRRset(g53.AnswerSection, rrsets[j])
		j++
	}
	for i := uint16(0); i < me.authorityCount; i++ {
		resp.AddRRset(g53.AuthSection, rrsets[j])
		j++
	}
	for i := uint16(0); i < me.additionalCount; i++ {
		resp.AddRRset(g53.AdditionalSection, rrsets[j])
		j++
	}
	resp.Header.Rcode = me.rcode
	resp.RecalculateSectionRRCount()
	return resp, true
}

func (c *MessageCache) get(keyHash, conflictHash uint64) (*MessageEntry, bool) {
	if elem, hit := c.data[keyHash]; hit {
		e := elem.Value.(*MessageEntry)
		if e.conflictHash == conflictHash && !e.IsExpire() {
			c.ll.MoveToFront(elem)
			return e, true
		}
	}
	return nil, false
}

func (c *MessageCache) Add(msg *g53.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if msg.Header.ANCount == 0 {
		e, res := negativeMessageToEntry(msg)
		c.add(e)
		c.rrsetCache.add(res)
	} else {
		e, res := positiveMsgToEntry(msg)
		c.add(e)
		c.rrsetCache.add(res)
	}
}

func positiveMsgToEntry(msg *g53.Message) (MessageEntry, []RRsetEntry) {
	keyHash, conflictHash := HashQuery(msg.Question.Name, msg.Question.Type)
	me := MessageEntry{
		keyHash:      keyHash,
		conflictHash: conflictHash,
		rcode:        msg.Header.Rcode,
	}

	rrCount := msg.Header.ANCount + msg.Header.NSCount + msg.Header.ARCount
	rrsets := make([]RRsetHash, 0, rrCount)
	rrsetEntries := make([]RRsetEntry, 0, rrCount)
	msgTtl := g53.RRTTL(math.MaxUint32)
	for _, sec := range []g53.SectionType{g53.AnswerSection, g53.AuthSection, g53.AdditionalSection} {
		for _, rrset := range msg.GetSection(sec) {
			if sec == g53.AnswerSection {
				me.answerCount += 1
			} else if sec == g53.AuthSection {
				me.authorityCount += 1
			} else {
				me.additionalCount += 1
			}

			keyHash, conflictHash := HashQuery(rrset.Name, rrset.Type)
			rrsets = append(rrsets, RRsetHash{
				keyHash:      keyHash,
				conflictHash: conflictHash,
			})

			if msgTtl > rrset.Ttl {
				msgTtl = rrset.Ttl
			}

			rrsetEntries = append(rrsetEntries, RRsetEntry{
				keyHash:      keyHash,
				conflictHash: conflictHash,
				rrset:        rrset,
				trustLevel:   getRRsetTrustLevel(msg, sec),
				expireTime:   time.Now().Add(time.Second * time.Duration(rrset.Ttl)),
			})
		}
	}
	me.rrsets = rrsets
	me.expireTime = time.Now().Add(time.Second * time.Duration(msgTtl))
	return me, rrsetEntries
}

func negativeMessageToEntry(msg *g53.Message) (MessageEntry, []RRsetEntry) {
	keyHash, conflictHash := HashQuery(msg.Question.Name, msg.Question.Type)
	me := MessageEntry{
		keyHash:         keyHash,
		conflictHash:    conflictHash,
		rcode:           msg.Header.Rcode,
		answerCount:     0,
		additionalCount: 0,
	}

	msgTtl := uint32(math.MaxUint32)
	auths := msg.GetSection(g53.AuthSection)
	rrsetEntries := make([]RRsetEntry, 0, 1)
	rrsets := make([]RRsetHash, 0, 1)
	if len(auths) == 1 && auths[0].Type == g53.RR_SOA && len(auths[0].Rdatas) == 1 {
		me.authorityCount = 1
		soa := auths[0]
		keyHash, conflictHash := HashQuery(soa.Name, soa.Type)
		rrsets = append(rrsets, RRsetHash{
			keyHash:      keyHash,
			conflictHash: conflictHash,
		})

		rrsetEntries = append(rrsetEntries, RRsetEntry{
			keyHash:      keyHash,
			conflictHash: conflictHash,
			rrset:        soa,
			trustLevel:   getRRsetTrustLevel(msg, g53.AuthSection),
			expireTime:   time.Now().Add(time.Second * time.Duration(soa.Ttl)),
		})

		rdata, _ := soa.Rdatas[0].(*g53.SOA)
		if msgTtl > rdata.Minimum {
			msgTtl = rdata.Minimum
		}
		if msgTtl > uint32(soa.Ttl) {
			msgTtl = uint32(soa.Ttl)
		}
	}
	me.rrsets = rrsets
	me.expireTime = time.Now().Add(time.Second * time.Duration(msgTtl))
	return me, rrsetEntries
}

func getRRsetTrustLevel(msg *g53.Message, sec g53.SectionType) TrustLevel {
	aa := msg.Header.GetFlag(g53.FLAG_AA)
	switch sec {
	case g53.AnswerSection:
		if aa {
			return AnswerWithAA
		} else {
			return AnswerWithoutAA
		}
	case g53.AuthSection:
		if aa {
			return AuthorityWithAA
		} else {
			return AuthorityWithoutAA
		}
	case g53.AdditionalSection:
		if aa {
			return AdditionalWithAA
		} else {
			return AdditionalWithoutAA
		}
	default:
		panic("unknown section type")
	}
}

func (c *MessageCache) add(e MessageEntry) {
	if elem, ok := c.data[e.keyHash]; ok {
		c.ll.MoveToFront(elem)
		elem.Value = &e
	} else if c.ll.Len() < c.cap {
		elem := c.ll.PushFront(&e)
		c.data[e.keyHash] = elem
	} else {
		//reuse last elem
		elem := c.ll.Back()
		oe := elem.Value.(*MessageEntry)
		delete(c.data, oe.keyHash)
		*oe = e
		c.data[e.keyHash] = elem
		c.ll.MoveToFront(elem)
	}
}

func (c *MessageCache) ResetCapacity(cap int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cap == cap {
		return
	}

	rc := c.ll.Len() - cap
	for rc > 0 {
		elem := c.ll.Back()
		oe := elem.Value.(*MessageEntry)
		c.remove(oe.keyHash, oe.conflictHash)
		rc -= 1
	}

	c.cap = cap
	c.rrsetCache.cap = RRsetVsMessageRatio * cap
}

func (c *MessageCache) Remove(name *g53.Name, typ g53.RRType) bool {
	keyHash, conflictHash := HashQuery(name, typ)
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.remove(keyHash, conflictHash)
}

func (c *MessageCache) remove(keyHash, conflictHash uint64) bool {
	if elem, hit := c.data[keyHash]; hit {
		e := elem.Value.(*MessageEntry)
		if e.conflictHash == conflictHash {
			delete(c.data, keyHash)
			c.ll.Remove(elem)
			return true
		}
	}
	return false
}
