package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type VLRHeader struct {
	Reserved                uint16
	UserID                  [16]byte
	RecordID                uint16
	RecordLengthAfterHeader uint16
	Description             [32]byte
}

func (vh *VLRHeader) getUserID() (userID string, err error) {
	chunks := bytes.Split(vh.UserID[:], []byte("\x00"))
	for _, chunk := range chunks {
		if len(chunk) != 0 {
			userID = string(chunk)
			return
		}
	}
	err = fmt.Errorf("userID is not set on VLR header")
	return
}

type VLR struct {
	header VLRHeader
	record []CRS
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
	bytesInRecord := make([]byte, v.header.RecordLengthAfterHeader)
	offsetToRecord := offsetIn + int64(binary.Size(VLRHeader{}))
	crs, err := v.getCRSFormat()
	if err != nil {
		return
	}
	_, err = file.ReadAt(bytesInRecord, offsetToRecord)
	if err != nil {
		return
	}
	crs.read(bytesInRecord, offsetToRecord)
	v.record = append(v.record, crs)
	offsetOut = offsetToRecord + int64(v.header.RecordLengthAfterHeader)
	return
}

func (v *VLR) getCRSFormat() (crs CRS, err error) {
	userID, err := v.header.getUserID()
	if err != nil {
		return
	}
	if userID == "LASF_Projection" {
		switch v.header.RecordID {
		case 34735:
			crs = &GeoKeyDirectoryTag{}
		case 34736:
			crs = &GeoDoubleParamsTag{}
		case 34737:
			crs = &GeoAsciiParamsTag{}
		case 2111:
			crs = &MathTransformWKT{}
		case 2112:
			crs = &CoordinateSystemWKT{}
		default:
			err = fmt.Errorf("CRS format not defined for LASF Projection")
		}
	} else if userID == "LASF_Spec" {
		switch v.header.RecordID {
		case 0:
			crs = &ClassificationLookup{}
		case 3:
			crs = &TextAreaDescription{}
		case 4:
			crs = &ExtraBytes{}
		default:
			err = fmt.Errorf("CRS format not defined for LASF Spec")
		}

	} else if userID == "liblas" {
		switch v.header.RecordID {
		case 2112:
			crs = &MathTransformWKT{}
		default:
			err = fmt.Errorf("CRS format not defined for liblas")
		}
	} else {
		err = fmt.Errorf("CRS format with userid %s not defined", userID)
	}
	return
}
