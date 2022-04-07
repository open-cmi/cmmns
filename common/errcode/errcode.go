package errcode

// 返回给前端的错误码
const (
	ErrFailed        = 1
	ErrNotRegistered = 1001
)

// 由xsnos返回的错误码
const (
	ERROR_CODE_SUCCESS     = 0
	ERROR_CODE_FAILED      = 0x20011001
	ERROR_CODE_ADDRESS_ERR = 0x20011002
	ERROR_CODE_CONN_ERR    = 0x20011003
	ERROR_CODE_WRITE_ERR   = 0x20011004
	ERROR_CODE_READ_ERR    = 0x20011005
)
