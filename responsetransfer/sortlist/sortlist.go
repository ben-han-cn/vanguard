package sortlist

import (
	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/vanguard/config"
	"github.com/ben-han-cn/vanguard/core"
)

type SortList struct {
	sorter RRsetSorter
}

func NewSortList() *SortList {
	return &SortList{
		sorter: newUserAddrBasedSorter(),
	}
}

func (m *SortList) ReloadConfig(conf *config.VanguardConf) {
	m.sorter.ReloadConfig(conf)
}

func (m *SortList) TransferResponse(client *core.Client) {
	if client.Response != nil {
		answers := client.Response.Sections[g53.AnswerSection]
		for i, rrset := range answers {
			answers[i] = m.sorter.Sort(client.View, client.IP(), rrset)
		}
	}
}
