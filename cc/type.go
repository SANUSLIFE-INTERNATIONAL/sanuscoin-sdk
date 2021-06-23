package cc

const (
	TypeMask     = 0xf0
	TransferMask = 0x10
	BurnMask     = 0x20
)

var (
	TransferOPCodes = [][]byte{
		[]byte{0x10}, // All Hashes in OP_RETURN}
		[]byte{0x11}, // SHA2 in Pay-to-Script-Hash multi-sig output (1 out of 2)
		[]byte{0x12}, // All Hashes in Pay-to-Script-Hash multi-sig outputs (1 out of 3)
		[]byte{0x13}, // Low security transaction no SHA2 for torrent data. SHA1 is always inside OP_RETURN in this case.
		[]byte{0x14}, // Low security transaction no SHA2 for torrent data. SHA1 is always inside OP_RETURN in this case. also no rules inside the metadata (if there are any they will be in ignored)
		[]byte{0x15}, // No metadata or rules (no SHA1 or SHA2)

	}

	BurnOPCodes = [][]byte{
		[]byte{0x20}, // All Hashes in OP_RETURN
		[]byte{0x21}, // SHA2 in Pay-to-Script-Hash multi-sig output (1 out of 2)
		[]byte{0x22}, // All Hashes in Pay-to-Script-Hash multi-sig outputs (1 out of 3)
		[]byte{0x23}, // Low security transaction no SHA2 for torrent data. SHA1 is always inside OP_RETURN in this case.
		[]byte{0x24}, // Low security transaction no SHA2 for torrent data. SHA1 is always inside OP_RETURN in this case. also no rules inside the metadata (if there are any they will be in ignored)
		[]byte{0x25}, // No metadata or rules (no SHA1 or SHA2)

	}
)
