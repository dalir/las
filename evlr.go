package las

import (
	"bytes"
	"encoding/binary"
	"os"
)

//  ______ __      __ _       _____
// |  ____|\ \    / /| |     |  __ \
// | |__    \ \  / / | |     | |__) |
// |  __|    \ \/ /  | |     |  _  /
// | |____    \  /   | |____ | | \ \
// |______|    \/    |______||_|  \_\
//
//

type EVLRHeader struct {
	Reserved                uint16
	UserID                  [16]byte
	RecordID                uint16
	RecordLengthAfterHeader uint64
	Description             [32]byte
}

type EVLR struct {
	header EVLRHeader
	record []byte
}

func (v *EVLR) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	headerInBytes := make([]byte, binary.Size(EVLRHeader{}))
	_, err = file.ReadAt(headerInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(headerInBytes), binary.LittleEndian, &v.header); err != nil {
		return
	}
	v.record = make([]byte, v.header.RecordLengthAfterHeader)
	offsetToRecord := offsetIn + int64(binary.Size(EVLRHeader{}))
	_, err = file.ReadAt(v.record, offsetToRecord)
	if err != nil {
		return
	}
	offsetOut = offsetToRecord + int64(v.header.RecordLengthAfterHeader)
	return
}
