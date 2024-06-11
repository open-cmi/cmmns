package wac

type SetRequest struct {
	Enable bool   `json:"enable"`
	Mode   string `json:"mode"`
}

type AddressRequest struct {
	Address string `json:"address"`
}
