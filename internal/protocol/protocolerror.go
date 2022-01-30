package protocol

const (
	StatusOK            = byte(0)
	StatusNotFound      = byte(1)
	StatusNotFit        = byte(2)
	StatusProtocolError = byte(3)
	StatusNoSpaceLeft   = byte(4)
	StatusUnknownError  = byte(255)
)
