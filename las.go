package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gocarina/gocsv"
	"math"
	"os"
)

const LAS_FILE_SIGNATURE = "LASF"

// _                _____
// | |        /\    / ____|
// | |       /  \  | (___
// | |      / /\ \  \___ \
// | |____ / ____ \ ____) |
// |______/_/    \_\_____/
//
//

type Las struct {
	header PublicHeaderBlock
	vlrs   []VLR
	pdrs   []PDR
	evlrs  []EVLR
}

func (l *Las) Parse(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	if err = l.readPublicHeaderBlock(file); err != nil {
		return
	}
	if err = l.readVLRs(file); err != nil {
		return
	}
	if err = l.readPDRs(file); err != nil {
		return
	}
	if err = l.readEVLRs(file); err != nil {
		return
	}
	return
}

func (l *Las) readPublicHeaderBlock(file *os.File) (err error) {
	headerInBytes := make([]byte, binary.Size(PublicHeaderBlock{}))
	_, err = file.ReadAt(headerInBytes, 0)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(headerInBytes), binary.LittleEndian, &l.header); err != nil {
		return
	}

	err = l.checkForCompliancy()
	if err != nil {
		return
	}

	return
}

func (l *Las) checkForCompliancy() (err error) {
	if err = l.isFileLasFormat(); err != nil {
		return
	}
	if err = l.isVersionOK(); err != nil {
		return
	}
	return
}

func (l *Las) isVersionOK() (err error) {
	fileVersion := l.header.GetVersion()
	for _, knownVersion := range AllLasVersions {
		if fileVersion == knownVersion {
			return
		}
	}
	err = fmt.Errorf("las files with version %s is not supported", fileVersion)
	return
}

func (l *Las) isFileLasFormat() (err error) {
	fileSignature := string(l.header.FileSignature[:])
	if fileSignature != LAS_FILE_SIGNATURE {
		err = fmt.Errorf("las files signature is not %s. File Signature: %s", LAS_FILE_SIGNATURE, fileSignature)

	}
	return
}

func (l *Las) readVLRs(file *os.File) (err error) {
	offset := int64(l.header.HeaderSize)
	for i := uint32(0); i < l.header.NumberOfVLRs; i++ {
		vlr := VLR{}
		offset, err = vlr.read(file, offset)
		if err != nil {
			return
		}
		l.vlrs = append(l.vlrs, vlr)
	}
	if uint32(offset) != l.header.OffsetToPointData {
		err = fmt.Errorf("after reading VLRs offset : %d, doesn't match the offset set in public header: %d", offset, l.header.OffsetToPointData)
	}
	return
}

func (l *Las) getPDRStruct() (pdr PDR) {
	switch l.header.PointDataRecordFormat {
	case 0:
		pdr = &PDR0{}
	case 1:
		pdr = &PDR1{}
	case 2:
		pdr = &PDR2{}
	case 3:
		pdr = &PDR3{}
	case 4:
		pdr = &PDR4{}
	case 5:
		pdr = &PDR5{}
	case 6:
		pdr = &PDR6{}
	case 7:
		pdr = &PDR7{}
	case 8:
		pdr = &PDR8{}
	case 9:
		pdr = &PDR9{}
	case 10:
		pdr = &PDR10{}
	}
	return
}

func (l *Las) getNumberOfPDRs() (numberOFPDRs uint64) {
	version := l.header.GetVersion()
	if version == V1_4 {
		numberOFPDRs = l.header.NumberOfPointRecords
	} else {
		numberOFPDRs = uint64(l.header.LegacyNumberOfPointRecords)
	}
	return
}

func (l *Las) readPDRs(file *os.File) (err error) {
	offset := int64(l.header.OffsetToPointData)
	numOfPDRs := l.getNumberOfPDRs()
	for i := uint64(0); i < numOfPDRs; i++ {
		pdr := l.getPDRStruct()
		offset, err = pdr.read(file, offset)
		if err != nil {
			return
		}
		l.pdrs = append(l.pdrs, pdr)
	}
	return
}

func (l *Las) getNumberOfEVLRs() (numberOFEVLRs uint32) {
	version := l.header.GetVersion()
	if version == V1_4 {
		numberOFEVLRs = l.header.NumberOfExtendedVariableLengthRecords
	} else {
		numberOFEVLRs = 0
	}
	return
}

func (l *Las) readEVLRs(file *os.File) (err error) {
	offset := int64(l.header.StartOfFirstExtendedVariableLengthRecord)
	numberOfEVLRs := l.getNumberOfEVLRs()
	for i := uint32(0); i < numberOfEVLRs; i++ {
		evlr := EVLR{}
		offset, err = evlr.read(file, offset)
		if err != nil {
			return
		}
		l.evlrs = append(l.evlrs, evlr)
	}
	return
}

func (l *Las) Las2txt(outputFile string) (err error) {
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	var csvList []XYZRGB

	for _, pdr := range l.pdrs {
		csvOutput := pdr.GetXYZRGB()

		csvOutput.X = math.Round((l.header.XOffset+csvOutput.X*l.header.XScaleFactor)*1000) / 1000
		csvOutput.Y = math.Round((l.header.YOffset+csvOutput.Y*l.header.YScaleFactor)*1000) / 1000
		csvOutput.Z = math.Round((l.header.ZOffset+csvOutput.Z*l.header.ZScaleFactor)*1000) / 1000

		csvOutput.R >>= 8
		csvOutput.G >>= 8
		csvOutput.B >>= 8

		csvList = append(csvList, csvOutput)
	}
	if err = gocsv.MarshalWithoutHeaders(csvList, file); err != nil {
		return
	}
	return
}
