package packet

type Packet struct {
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

	data[ISEQ] = p.SEQ
	data[ITYP] = p.TYP
	data[ILEN] = p.LEN
	data[IRET] = p.RET
	data[ITMO] = p.TMO

	if p.LEN > 0 {
		data = append(data, p.PYL...)
	}

	return data
}

func Decode(data []byte) Packet {
	p := Packet{}

	p.SEQ = data[ISEQ]
	p.TYP = data[ITYP]
	p.LEN = data[ILEN]
	p.RET = data[IRET]
	p.TMO = data[ITMO]

	if p.LEN > 0 {
		p.PYL = data[IPYL:]
	}

	return p
}
