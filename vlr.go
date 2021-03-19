package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

type VLRHeader struct {
	Reserved                uint16
	UserID                  [16]byte
	RecordID                uint16
	RecordLengthAfterHeader uint16
	Description             [32]byte
}

type VLR struct {
	header VLRHeader
	record []byte
}

func (v *VLR) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	headerInBytes := make([]byte, binary.Size(VLRHeader{}))
	_, err = file.ReadAt(headerInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(headerInBytes), binary.LittleEndian, &v.header); err != nil {
		return
	}
	v.record = make([]byte, v.header.RecordLengthAfterHeader)
	offsetToRecord := offsetIn + int64(binary.Size(VLRHeader{}))
	_, err = file.ReadAt(v.record, offsetToRecord)
	if err != nil {
		return
	}
	offsetOut = offsetToRecord + int64(v.header.RecordLengthAfterHeader)
	return
}
