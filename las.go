package las

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
	Header PublicHeaderBlock
	Vlrs   []VLR
	Pdrs   PDRs
	Evlrs  []EVLR
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
	if err = binary.Read(bytes.NewReader(headerInBytes), binary.LittleEndian, &l.Header); err != nil {
		return
	}

	err = l.checkForCompliancy()
	if err != nil {
		return
	}

	return
}
func (l *Las) WritePHB(file *os.File) (err error) {
	headerInBytes := make([]byte, binary.Size(PublicHeaderBlock{}))
	_, err = file.WriteAt(headerInBytes, 0)
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
	fileVersion := l.Header.GetVersion()
	for _, knownVersion := range AllLasVersions {
		if fileVersion == knownVersion {
			return
		}
	}
	err = fmt.Errorf("las files with version %s is not supported", fileVersion)
	return
}

func (l *Las) isFileCompressed() (compressed bool) {
	compressed = false
	if l.Header.PointDataRecordFormat&0x80 == 0x80 {
		compressed = true
	}
	return
}

func (l *Las) isFileLasFormat() (err error) {
	fileSignature := string(l.Header.FileSignature[:])
	if fileSignature != LAS_FILE_SIGNATURE {
		err = fmt.Errorf("las files signature is not %s. File Signature: %s", LAS_FILE_SIGNATURE, fileSignature)
	}
	if l.isFileCompressed() {
		err = fmt.Errorf("file is compressed. Needs to be treated as laz format")
	}
	return
}

func (l *Las) readVLRs(file *os.File) (err error) {
	offset := int64(l.Header.HeaderSize)
	for i := uint32(0); i < l.Header.NumberOfVLRs; i++ {
		vlr := VLR{}
		offset, err = vlr.read(file, offset)
		if err != nil {
			return
		}
		l.Vlrs = append(l.Vlrs, vlr)
	}
	if uint32(offset) != l.Header.OffsetToPointData {
		err = fmt.Errorf("after reading VLRs offset : %d, doesn't match the offset set in public Header: %d", offset, l.Header.OffsetToPointData)
	}
	return
}

func (l *Las) getNumberOfPDRs() (numberOFPDRs uint64) {
	version := l.Header.GetVersion()
	if version == V1_4 {
		numberOFPDRs = l.Header.NumberOfPointRecords
	} else {
		numberOFPDRs = uint64(l.Header.LegacyNumberOfPointRecords)
	}
	return
}

func (l *Las) readPDRs(file *os.File) (err error) {
	numOfPDRs := l.getNumberOfPDRs()
	switch l.Header.PointDataRecordFormat {
	case 0:
		l.Pdrs = make(PDR0s, numOfPDRs)
	case 1:
		l.Pdrs = make(PDR1s, numOfPDRs)
	case 2:
		l.Pdrs = make(PDR2s, numOfPDRs)
	case 3:
		l.Pdrs = make(PDR3s, numOfPDRs)
	case 4:
		l.Pdrs = make(PDR4s, numOfPDRs)
	case 5:
		l.Pdrs = make(PDR5s, numOfPDRs)
	case 6:
		l.Pdrs = make(PDR6s, numOfPDRs)
	case 7:
		l.Pdrs = make(PDR7s, numOfPDRs)
	case 8:
		l.Pdrs = make(PDR8s, numOfPDRs)
	case 9:
		l.Pdrs = make(PDR9s, numOfPDRs)
	case 10:
		l.Pdrs = make(PDR10s, numOfPDRs)
	default:
		err = fmt.Errorf("point data record format not recognised")
		return
	}
	if err = l.Pdrs.read(file, int64(l.Header.OffsetToPointData), uint64(l.Header.PointDataRecordLength)); err != nil {
		return
	}
	return
}

func (l *Las) getNumberOfEVLRs() (numberOFEVLRs uint32) {
	version := l.Header.GetVersion()
	if version == V1_4 {
		numberOFEVLRs = l.Header.NumberOfExtendedVariableLengthRecords
	} else {
		numberOFEVLRs = 0
	}
	return
}

func (l *Las) readEVLRs(file *os.File) (err error) {
	offset := int64(l.Header.StartOfFirstExtendedVariableLengthRecord)
	numberOfEVLRs := l.getNumberOfEVLRs()
	for i := uint32(0); i < numberOfEVLRs; i++ {
		evlr := EVLR{}
		offset, err = evlr.read(file, offset)
		if err != nil {
			return
		}
		l.Evlrs = append(l.Evlrs, evlr)
	}
	return
}

func (l *Las) Las2txt(outputFile string) (err error) {
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	csvList := l.Pdrs.GetCSVList()

	for _, csvOutput := range csvList {

		csvOutput.X = math.Round((l.Header.XOffset+csvOutput.X*l.Header.XScaleFactor)*1000) / 1000
		csvOutput.Y = math.Round((l.Header.YOffset+csvOutput.Y*l.Header.YScaleFactor)*1000) / 1000
		csvOutput.Z = math.Round((l.Header.ZOffset+csvOutput.Z*l.Header.ZScaleFactor)*1000) / 1000

		if (csvOutput.R|csvOutput.G|csvOutput.B)&0xFF00 == csvOutput.R|csvOutput.G|csvOutput.B { // Checking the 8 bit channel vs 16 bit
			csvOutput.R >>= 8
			csvOutput.G >>= 8
			csvOutput.B >>= 8
		}
	}

	if err = gocsv.MarshalWithoutHeaders(csvList, file); err != nil {
		return
	}
	return
}
