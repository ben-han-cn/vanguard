package localdata

import (
	"github.com/ben-han-cn/g53"
	"github.com/ben-han-cn/vanguard/core"
)

func realQueryName(client *core.Client) *g53.Name {
	if client.Response == nil {
		return client.Request.Question.Name
	}

	answer := client.Response.Sections[g53.AnswerSection]
	answerCount := len(answer)
	if answerCount == 0 || answer[answerCount-1].Type != g53.RR_CNAME {
		return client.Request.Question.Name
	}

	return answer[answerCount-1].Rdatas[0].(*g53.CName).Name
}
