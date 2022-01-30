package protocol

const (
	Ping   = byte(0)
	Create = byte(1)
	Update = byte(2)
	Delete = byte(3)
	Read   = byte(4)

	defaultBuffSize = 65532
)
