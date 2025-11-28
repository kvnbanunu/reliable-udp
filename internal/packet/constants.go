package packet

import "errors"

const (
	SYN uint8 = iota // 0
	ACK              // 1
	SND              // 2
	FIN              // 3
)

const (
	HEADER_LEN      int = 6
	MAX_PAYLOAD_LEN int = 255
	MAX_PACKET_LEN  int = 261 // 5 byte header + 255 byte payload
)

// Indexes for packet fields
const (
	ICID int = iota // 0
	ISEQ            // 1
	ITYP            // 2
	ILEN            // 3
	IRET            // 4
	ITMO            // 5
	IPYL            // 6
)

var (
	ErrTimeout = errors.New("Error: Time out")
	ErrNoRead  = errors.New("Error: No bytes read")
	ErrNoWrite = errors.New("Error: No bytes written")
	ErrLongMsg = errors.New("Error: Input message is too long")
)
