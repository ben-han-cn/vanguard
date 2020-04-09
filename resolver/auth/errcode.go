package auth

import (
	"github.com/ben-han-cn/vanguard/httpcmd"
)

var (
	ErrNonExistZone        = httpcmd.NewError(httpcmd.AuthErrCodeStart, "operate non-exist zone")
	ErrGetZoneFail         = httpcmd.NewError(httpcmd.AuthErrCodeStart+1, "get zone from backend server failed")
	ErrZoneExists          = httpcmd.NewError(httpcmd.AuthErrCodeStart+2, "zone already exists")
	ErrUnSupportMetaType   = httpcmd.NewError(httpcmd.AuthErrCodeStart+3, "unsupported meta data type")
	ErrNonExistRR          = httpcmd.NewError(httpcmd.AuthErrCodeStart+4, "operate non-exist rr")
	ErrUnSupportZoneType   = httpcmd.NewError(httpcmd.AuthErrCodeStart+5, "not supported zone type")
	ErrBadZoneViewOwner    = httpcmd.NewError(httpcmd.AuthErrCodeStart+6, "zone owner doesn't has view owner")
	ErrZoneUpdateFailed    = httpcmd.NewError(httpcmd.AuthErrCodeStart+7, "update zone failed")
	ErrAddDuplicateRR      = httpcmd.NewError(httpcmd.AuthErrCodeStart+8, "add duplicate rr")
	ErrAddExcludeRR        = httpcmd.NewError(httpcmd.AuthErrCodeStart+9, "add exclusive rr")
	ErrShortOfGlue         = httpcmd.NewError(httpcmd.AuthErrCodeStart+10, "short of glue rr")
	ErrAlreadyHasCNAME     = httpcmd.NewError(httpcmd.AuthErrCodeStart+11, "conflict with exists cname")
	ErrDeleteUnknownRR     = httpcmd.NewError(httpcmd.AuthErrCodeStart+12, "delete unknown rr")
	ErrDeleteSOARecord     = httpcmd.NewError(httpcmd.AuthErrCodeStart+13, "can't delete soa rr")
	ErrDeleteOnlyLeftNS    = httpcmd.NewError(httpcmd.AuthErrCodeStart+14, "no ns left after delete")
	ErrDeleteNecessaryGlue = httpcmd.NewError(httpcmd.AuthErrCodeStart+15, "delete glue needed by other rr")
	ErrNonExistReverseZone = httpcmd.NewError(httpcmd.AuthErrCodeStart+16, "reverse zone doesn't exist")
	ErrRdataFormatError    = httpcmd.NewError(httpcmd.AuthErrCodeStart+17, "rdata is invalid")
	ErrAddOutofZoneRR      = httpcmd.NewError(httpcmd.AuthErrCodeStart+18, "rr is out of zone")
	ErrUnspportRRType      = httpcmd.NewError(httpcmd.AuthErrCodeStart+19, "rr type isn't supported")
	ErrUpdateSlaveZone     = httpcmd.NewError(httpcmd.AuthErrCodeStart+20, "can't update slave zone")
	ErrInvalidZoneName     = httpcmd.NewError(httpcmd.AuthErrCodeStart+21, "zone name isn't valid")
	ErrSOASerialDegrade    = httpcmd.NewError(httpcmd.AuthErrCodeStart+22, "soa serial number degraded")
	ErrAuthZoneExists      = httpcmd.NewError(httpcmd.AuthErrCodeStart+23, "auth zone with same name already exists")
	ErrInvalidZoneData     = httpcmd.NewError(httpcmd.AuthErrCodeStart+24, "invalid zone data")
	ErrInvalidRR           = httpcmd.NewError(httpcmd.AuthErrCodeStart+25, "rr data isn't valid")
	ErrDeleteZoneFailed    = httpcmd.NewError(httpcmd.AuthErrCodeStart+26, "delete auth zone failed")
	ErrUpdateZoneFailed    = httpcmd.NewError(httpcmd.AuthErrCodeStart+27, "update auth zone failed")
)
