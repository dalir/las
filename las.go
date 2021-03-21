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
	pdrs   PDRs
	evlrs  []EVLR
}

func (l *Las) Parse(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	if err = l.readPHB(file); err != nil {
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

func (l *Las) readPHB(file *os.File) (err error) {
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
	numOfPDRs := l.getNumberOfPDRs()
	switch l.header.PointDataRecordFormat {
	case 0:
		l.pdrs = make(PDR0s, numOfPDRs)
	case 1:
		l.pdrs = make(PDR1s, numOfPDRs)
	case 2:
		l.pdrs = make(PDR2s, numOfPDRs)
	case 3:
		l.pdrs = make(PDR3s, numOfPDRs)
	case 4:
		l.pdrs = make(PDR4s, numOfPDRs)
	case 5:
		l.pdrs = make(PDR5s, numOfPDRs)
	case 6:
		l.pdrs = make(PDR6s, numOfPDRs)
	case 7:
		l.pdrs = make(PDR7s, numOfPDRs)
	case 8:
		l.pdrs = make(PDR8s, numOfPDRs)
	case 9:
		l.pdrs = make(PDR9s, numOfPDRs)
	case 10:
		l.pdrs = make(PDR10s, numOfPDRs)
	default:
		err = fmt.Errorf("point data record format not recognised")
		return
	}
	if err = l.pdrs.read(file, int64(l.header.OffsetToPointData)); err != nil {
		return
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

	csvList := l.pdrs.GetCSVList()

	for _, csvOutput := range csvList {

		csvOutput.X = math.Round((l.header.XOffset+csvOutput.X*l.header.XScaleFactor)*1000) / 1000
		csvOutput.Y = math.Round((l.header.YOffset+csvOutput.Y*l.header.YScaleFactor)*1000) / 1000
		csvOutput.Z = math.Round((l.header.ZOffset+csvOutput.Z*l.header.ZScaleFactor)*1000) / 1000

		csvOutput.R >>= 8
		csvOutput.G >>= 8
		csvOutput.B >>= 8
	}

	if err = gocsv.MarshalWithoutHeaders(csvList, file); err != nil {
		return
	}
	return
}
