package packet

import "errors"

const (
	SND uint8 = iota // 0
	ACK              // 1
	FIN              // 2
)

const (
	HEADER_LEN      int = 5
	MAX_PAYLOAD_LEN int = 255
	MAX_PACKET_LEN  int = 260 // 5 byte header + 255 byte payload
)

// Indexes for packet fields
const (
	ISEQ int = iota // 0
	ITYP            // 1
	ILEN            // 2
	IRET            // 3
	ITMO            // 4
	IPYL            // 5
)

var (
	ErrBadReq  = errors.New("Bad Request")
	ErrTimeout = errors.New("Time out")
	ErrNoRead  = errors.New("No bytes read")
	ErrNoWrite = errors.New("No bytes written")
	ErrLongMsg = errors.New("Input message is too long")
	ErrInvTYP  = errors.New("Invalid Packet Type")
	ErrDupPCK  = errors.New("Duplicate or old packet")
	ErrCancel  = errors.New("Message canceled")
)
