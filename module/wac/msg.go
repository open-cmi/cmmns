package wac

type SetRequest struct {
	Mode             string `json:"mode"`
	RawPermitAddress string `json:"raw_permit_address"`
	RawDenyAddress   string `json:"deny_permit_address"`
}
