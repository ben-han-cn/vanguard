package cache

import (
	"testing"
	"time"

	ut "github.com/ben-han-cn/cement/unittest"
	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/vanguard/logger"
)

func buildMessage(name_ string, ips []string, ttl int) *g53.Message {
	header := g53.Header{
		Id:      1000,
		Opcode:  g53.OP_QUERY,
		Rcode:   g53.R_NOERROR,
		QDCount: 1,
		ANCount: 2,
		NSCount: 0,
		ARCount: 1,
	}

	name := g53.NameFromStringUnsafe(name_)
	question := &g53.Question{
		Name:  name,
		Type:  g53.RR_A,
		Class: g53.CLASS_IN,
	}

	rdatas := make([]g53.Rdata, 0, len(ips))
	for _, ip := range ips {
		rdata, _ := g53.AFromString(ip)
		rdatas = append(rdatas, rdata)
	}

	var answer g53.Section
	answer = append(answer, &g53.RRset{
		Name:   name,
		Type:   g53.RR_A,
		Class:  g53.CLASS_IN,
		Ttl:    g53.RRTTL(ttl),
		Rdatas: rdatas,
	})

	return &g53.Message{
		Header:   header,
		Question: question,
		Sections: [...]g53.Section{answer, []*g53.RRset{}, []*g53.RRset{}},
		Edns: &g53.EDNS{
			Version:     0,
			UdpSize:     4096,
			DnssecAware: false,
		},
	}
}

func TestMessageCache(t *testing.T) {
	logger.UseDefaultLogger("debug")
	cache := newMessageCache(3)

	ut.Equal(t, cache.Len(), 0)

	message := buildMessage("test.example.com.", []string{"1.1.1.1"}, 3)
	cache.Add(message)
	ut.Equal(t, cache.Len(), 1)

	qname, _ := g53.NameFromString("test.example.com.")
	request := g53.MakeQuery(qname, g53.RR_A, 512, false)
	request.Header.Id = 1000
	message, found := cache.Get(request)
	ut.Assert(t, found == true, "message should be fetched")
	ut.Equal(t, message.Header.Id, uint16(1000))

	cache.Add(message)
	ut.Equal(t, cache.Len(), 1)

	message1 := buildMessage("test1.example.com.", []string{"1.1.1.1"}, 3)
	cache.Add(message1)
	ut.Equal(t, cache.Len(), 2)
	message2 := buildMessage("test2.example.com.", []string{"1.1.1.1"}, 3)
	cache.Add(message2)
	ut.Equal(t, cache.Len(), 3)

	message3 := buildMessage("test3.example.com.", []string{"1.1.1.1"}, 3)
	cache.Add(message3)
	ut.Equal(t, cache.Len(), 3)

	<-time.After(4 * time.Second)
	_, found = cache.Get(request)
	ut.Assert(t, found == false, "message should expired")
	ut.Equal(t, cache.Len(), 3)

	cache.Add(buildMessage("test.example.com.", []string{"2.2.2.2", "1.1.1.1"}, 30))
	ut.Equal(t, cache.Len(), 3)
	message, found = cache.Get(request)
	ut.Assert(t, found == true, "message shouldn't expired")
	ut.Equal(t, message.Sections[g53.AnswerSection][0].Rdatas[0].String(), "2.2.2.2")

	cache.Remove(qname, g53.RR_A)
	_, found = cache.Get(request)
	ut.Assert(t, found == false, "message should be cleaned")
	ut.Equal(t, cache.Len(), 2)
}
