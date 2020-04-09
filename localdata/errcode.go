package localdata

import (
	"github.com/ben-han-cn/vanguard/httpcmd"
)

var (
	ErrInvalidPolicy   = httpcmd.NewError(httpcmd.LocalDataErrCodeStart, "policy data isn't valid")
	ErrDuplicatePolicy = httpcmd.NewError(httpcmd.LocalDataErrCodeStart+1, "duplicate policy")
)
