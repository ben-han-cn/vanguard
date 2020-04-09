package querysource

import (
	"github.com/ben-han-cn/vanguard/httpcmd"
)

var (
	ErrDuplicateQuerySource = httpcmd.NewError(httpcmd.QuerySourceErrCodeStart, "duplicate query source")
	ErrNotExistQuerySource  = httpcmd.NewError(httpcmd.QuerySourceErrCodeStart+1, "unknown query source")
)
