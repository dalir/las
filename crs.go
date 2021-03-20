package main

import (
	"bytes"
	"encoding/binary"
)

type CRS interface {
	read(record []byte) (err error)
}

type GeoKeyDirectoryTag struct {
	header GeoKeyDirectoryHeader
	keys   []KeyEntry
}

type GeoKeyDirectoryHeader struct {
	KeyDirectoryVersion uint16
	KeyRevision         uint16
	MinorRevision       uint16
	NumberOfKeys        uint16
}

type KeyEntry struct {
	KeyID           uint16
	TIFFTagLocation uint16
	Count           uint16
	ValueOffset     uint16
}

func (g *GeoKeyDirectoryTag) read(record []byte) (err error) {
	headerSize := binary.Size(GeoKeyDirectoryHeader{})
	if err = binary.Read(bytes.NewReader(record[:headerSize]), binary.LittleEndian, &g.header); err != nil {
		return
	}
	g.keys = make([]KeyEntry, g.header.NumberOfKeys)
	if err = binary.Read(bytes.NewReader(record[headerSize:]), binary.LittleEndian, &g.keys); err != nil {
		return
	}
	return
}
