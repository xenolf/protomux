package protomux

type tlsProbe struct {
}

func (p tlsProbe) apply(data []byte) bool {
	// TLS packet starts with a record "Hello" (0x16), followed by version
	// (0x03 0x00-0x03) (RFC6101 A.1)
	// This means we reject SSLv2 and lower, which is actually a good thing (RFC6176)
	return data[0] == 0x16 && data[1] == 0x03 && (data[2] >= 0 && data[2] <= 0x03)
}
