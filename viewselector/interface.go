package viewselector

import (
	"github.com/ben-han-cn/vanguard/config"
	"github.com/ben-han-cn/vanguard/core"
)

type ViewSelector interface {
	ReloadConfig(*config.VanguardConf)
	ViewForQuery(*core.Client) (string, bool)
	GetViews() []string
}
