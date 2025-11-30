package packet

// Packet Type
type PType uint8

// Header for this new protocol
// Total size is 5 bytes
type Packet struct {
	CID uint8  // Client ID, assigned by the server | default 0 for new connections
	SEQ uint8  // Sequence Number
	TYP uint8  // Packet Type
	LEN uint8  // Length of the payload | Max 255
	RET uint8  // Number of resends attempted
	TMO uint8  // Timeout in seconds
	PYL []byte // Payload max 255 bytes
}

func NewPacket(seq, typ, tmo uint8, msg string) (Packet, error) {
	p := Packet{}
	length := len(msg)

	if length > MAX_PAYLOAD_LEN {
		return p, ErrLongMsg
	}

	p.CID = 0
	p.SEQ = seq
	p.TYP = typ
	p.LEN = uint8(length)
	p.RET = 0
	p.TMO = tmo
	if length > 0 {
		p.PYL = []byte(msg)
	}

	return p, nil
}

func Encode(p Packet) []byte {
	length := HEADER_LEN + int(p.LEN)
	data := make([]byte, length)

	data[ICID] = p.CID
	data[ISEQ] = p.SEQ
	data[ITYP] = p.TYP
	data[ILEN] = p.LEN
	data[IRET] = p.RET
	data[ITMO] = p.TMO

	if p.TYP == SND && p.LEN > 0 {
		data = append(data, p.PYL...)
	}

	return data
}

func Decode(data []byte) Packet {
	p := Packet{}

	p.CID = data[ICID]
	p.SEQ = data[ISEQ]
	p.TYP = data[ITYP]
	p.LEN = data[ILEN]
	p.RET = data[IRET]
	p.TMO = data[ITMO]

	if p.TYP == SND && p.LEN > 0 {
		p.PYL = data[IPYL:]
	}

	return p
}
