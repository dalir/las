package main

import (
	"bytes"
	"encoding/binary"
)

//   _____ _____   _____
//  / ____|  __ \ / ____|
// | |    | |__) | (___
// | |    |  _  / \___ \
// | |____| | \ \ ____) |
//  \_____|_|  \_\_____/
//
//

type CRS interface {
	read(record []byte, offset int64) (err error)
}

//   _____            _  __          _____  _               _                _______
//  / ____|          | |/ /         |  __ \(_)             | |              |__   __|
// | |  __  ___  ___ | ' / ___ _   _| |  | |_ _ __ ___  ___| |_ ___  _ __ _   _| | __ _  __ _
// | | |_ |/ _ \/ _ \|  < / _ \ | | | |  | | | '__/ _ \/ __| __/ _ \| '__| | | | |/ _` |/ _` |
// | |__| |  __/ (_) | . \  __/ |_| | |__| | | | |  __/ (__| || (_) | |  | |_| | | (_| | (_| |
//  \_____|\___|\___/|_|\_\___|\__, |_____/|_|_|  \___|\___|\__\___/|_|   \__, |_|\__,_|\__, |
//                              __/ |                                      __/ |         __/ |
//                             |___/                                      |___/         |___/

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

func (g *GeoKeyDirectoryTag) read(record []byte, offset int64) (err error) {
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

//   _____            _____              _     _      _____                           _______
//  / ____|          |  __ \            | |   | |    |  __ \                         |__   __|
// | |  __  ___  ___ | |  | | ___  _   _| |__ | | ___| |__) |_ _ _ __ __ _ _ __ ___  ___| | __ _  __ _
// | | |_ |/ _ \/ _ \| |  | |/ _ \| | | | '_ \| |/ _ \  ___/ _` | '__/ _` | '_ ` _ \/ __| |/ _` |/ _` |
// | |__| |  __/ (_) | |__| | (_) | |_| | |_) | |  __/ |  | (_| | | | (_| | | | | | \__ \ | (_| | (_| |
//  \_____|\___|\___/|_____/ \___/ \__,_|_.__/|_|\___|_|   \__,_|_|  \__,_|_| |_| |_|___/_|\__,_|\__, |
//                                                                                                __/ |
//                                                                                               |___/

type GeoDoubleParamsTag map[int64]float64

func (g GeoDoubleParamsTag) read(record []byte, offset int64) (err error) {
	for recordLocation := offset; recordLocation < int64(len(record))+offset; recordLocation += 4 {
		var value float64
		if err = binary.Read(bytes.NewReader(record[recordLocation:recordLocation+4]), binary.LittleEndian, &value); err != nil {
			return
		}
		g[recordLocation] = value
	}
	return
}

//   _____                             _ _ _____                           _______
//  / ____|             /\            (_|_)  __ \                         |__   __|
// | |  __  ___  ___   /  \   ___  ___ _ _| |__) |_ _ _ __ __ _ _ __ ___  ___| | __ _  __ _
// | | |_ |/ _ \/ _ \ / /\ \ / __|/ __| | |  ___/ _` | '__/ _` | '_ ` _ \/ __| |/ _` |/ _` |
// | |__| |  __/ (_) / ____ \\__ \ (__| | | |  | (_| | | | (_| | | | | | \__ \ | (_| | (_| |
//  \_____|\___|\___/_/    \_\___/\___|_|_|_|   \__,_|_|  \__,_|_| |_| |_|___/_|\__,_|\__, |
//                                                                                     __/ |
//                                                                                    |___/

type GeoAsciiParamsTag map[int64]string

func (g GeoAsciiParamsTag) read(record []byte, offset int64) (err error) {
	chunks := bytes.Split(record, []byte("\x00"))
	for _, chunk := range chunks {
		if len(chunk) != 0 {
			g[offset] = string(chunk)
			offset += int64(len(chunk) + 1)
		}
	}
	return
}

//  __  __       _   _  _______                   __                 __          ___  _________
// |  \/  |     | | | ||__   __|                 / _|                \ \        / / |/ /__   __|
// | \  / | __ _| |_| |__ | |_ __ __ _ _ __  ___| |_ ___  _ __ _ __ __\ \  /\  / /| ' /   | |
// | |\/| |/ _` | __| '_ \| | '__/ _` | '_ \/ __|  _/ _ \| '__| '_ ` _ \ \/  \/ / |  <    | |
// | |  | | (_| | |_| | | | | | | (_| | | | \__ \ || (_) | |  | | | | | \  /\  /  | . \   | |
// |_|  |_|\__,_|\__|_| |_|_|_|  \__,_|_| |_|___/_| \___/|_|  |_| |_| |_|\/  \/   |_|\_\  |_|
//
//

type MathTransformWKT []string

func (m *MathTransformWKT) read(record []byte, offset int64) (err error) {
	chunks := bytes.Split(record, []byte("\x00"))
	for _, chunk := range chunks {
		if len(chunk) != 0 {
			*m = append(*m, string(chunk))
		}
	}
	return
}

//   _____                    _ _             _        _____           _             __          ___  _________
//  / ____|                  | (_)           | |      / ____|         | |            \ \        / / |/ /__   __|
// | |     ___   ___  _ __ __| |_ _ __   __ _| |_ ___| (___  _   _ ___| |_ ___ _ __ __\ \  /\  / /| ' /   | |
// | |    / _ \ / _ \| '__/ _` | | '_ \ / _` | __/ _ \\___ \| | | / __| __/ _ \ '_ ` _ \ \/  \/ / |  <    | |
// | |___| (_) | (_) | | | (_| | | | | | (_| | ||  __/____) | |_| \__ \ ||  __/ | | | | \  /\  /  | . \   | |
//  \_____\___/ \___/|_|  \__,_|_|_| |_|\__,_|\__\___|_____/ \__, |___/\__\___|_| |_| |_|\/  \/   |_|\_\  |_|
//                                                            __/ |
//                                                           |___/

type CoordinateSystemWKT []string

func (c *CoordinateSystemWKT) read(record []byte, offset int64) (err error) {
	chunks := bytes.Split(record, []byte("\x00"))
	for _, chunk := range chunks {
		if len(chunk) != 0 {
			*c = append(*c, string(chunk))
		}
	}
	return
}

//   _____ _               _  __ _           _   _               _                 _
//  / ____| |             (_)/ _(_)         | | (_)             | |               | |
// | |    | | __ _ ___ ___ _| |_ _  ___ __ _| |_ _  ___  _ __   | |     ___   ___ | | ___   _ _ __
// | |    | |/ _` / __/ __| |  _| |/ __/ _` | __| |/ _ \| '_ \  | |    / _ \ / _ \| |/ / | | | '_ \
// | |____| | (_| \__ \__ \ | | | | (_| (_| | |_| | (_) | | | | | |___| (_) | (_) |   <| |_| | |_) |
//  \_____|_|\__,_|___/___/_|_| |_|\___\__,_|\__|_|\___/|_| |_| |______\___/ \___/|_|\_\\__,_| .__/
//                                                                                           | |
//                                                                                           |_|

type ClassificationLookup [256]Classification

type Classification struct {
	ClassNumber uint8
	Description [15]byte
}

func (c ClassificationLookup) read(record []byte, offset int64) (err error) {
	if err = binary.Read(bytes.NewReader(record), binary.LittleEndian, &c); err != nil {
		return
	}
	return
}

//  _______        _                               _____                      _       _   _
// |__   __|      | |       /\                    |  __ \                    (_)     | | (_)
//    | | _____  _| |_     /  \   _ __ ___  __ _  | |  | | ___  ___  ___ _ __ _ _ __ | |_ _  ___  _ __
//    | |/ _ \ \/ / __|   / /\ \ | '__/ _ \/ _` | | |  | |/ _ \/ __|/ __| '__| | '_ \| __| |/ _ \| '_ \
//    | |  __/>  <| |_   / ____ \| | |  __/ (_| | | |__| |  __/\__ \ (__| |  | | |_) | |_| | (_) | | | |
//    |_|\___/_/\_\\__| /_/    \_\_|  \___|\__,_| |_____/ \___||___/\___|_|  |_| .__/ \__|_|\___/|_| |_|
//                                                                             | |
//                                                                             |_|

type TextAreaDescription []string

func (t *TextAreaDescription) read(record []byte, offset int64) (err error) {
	chunks := bytes.Split(record, []byte("\x00"))
	for _, chunk := range chunks {
		if len(chunk) != 0 {
			*t = append(*t, string(chunk))
		}
	}
	return
}

//  ______      _               ____        _
// |  ____|    | |             |  _ \      | |
// | |__  __  _| |_ _ __ __ _  | |_) |_   _| |_ ___  ___
// |  __| \ \/ / __| '__/ _` | |  _ <| | | | __/ _ \/ __|
// | |____ >  <| |_| | | (_| | | |_) | |_| | ||  __/\__ \
// |______/_/\_\\__|_|  \__,_| |____/ \__, |\__\___||___/
//                                     __/ |
//                                    |___/

type ExtraBytes []ExtraBytesDescriptor

type ExtraBytesDescriptor struct {
	reserved    [2]uint8
	dataType    uint8
	options     uint8
	name        [32]byte
	unused      [4]uint8
	noData      [8]byte
	deprecated1 [16]uint8
	min         [8]byte
	deprecated2 [16]uint8
	max         [8]byte
	deprecated3 [16]uint8
	scale       float64
	deprecated4 [16]uint8
	offset      float64
	deprecated5 [16]uint8
	description [32]byte
}

func (e *ExtraBytes) read(record []byte, offset int64) (err error) {
	sizeOfDescriptor := int64(binary.Size(ExtraBytesDescriptor{}))
	for recordLocation := offset; recordLocation < int64(len(record))+offset; recordLocation += sizeOfDescriptor {
		value := ExtraBytesDescriptor{}
		if err = binary.Read(bytes.NewReader(record[recordLocation:recordLocation+sizeOfDescriptor]), binary.LittleEndian, &value); err != nil {
			return
		}
		*e = append(*e, value)
	}
	return
}
