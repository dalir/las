package las

import (
	"bytes"
	"encoding/binary"
	"os"
)

// _____  _____  _____
// |  __ \|  __ \|  __ \
// | |__) | |  | | |__) |
// |  ___/| |  | |  _  /
// | |    | |__| | | \ \
// |_|    |_____/|_|  \_\
//
//

type PDRs interface {
	read(file *os.File, offsetIn int64, dataLength uint64) (err error)
	GetCSVList() (output []*XYZRGB)
}

type XYZRGB struct {
	X float64 `csv:"X"`
	Y float64 `csv:"Y"`
	Z float64 `csv:"Z"`
	R uint16  `csv:"R"`
	G uint16  `csv:"G"`
	B uint16  `csv:"B"`
}

//  ______  ____   _____   __  __         _______    ___
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __|  / _ \
// | |__  | |  | || |__) || \  / |   /  \   | |    | | | |
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |    | | | |
// | |    | |__| || | \ \ | |  | | / ____ \ | |    | |_| |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|     \___/
//
//

const (
	PDR0_RETURN_NUMBER_MASK       = 0x07
	PDR0_NUMBER_OF_RETURNS_MASK   = 0x38
	PDR0_SCAN_DIRECTION_FLAG_MASK = 0x40
	PDR0_EDGE_OF_FLIGHT_LINE_MASK = 0x80
)

type ClassAttribute uint8

const (
	Created_NeverClassified ClassAttribute = iota
	Uncalassified
	Ground
	Low_Vegetation
	Medium_Vegetation
	High_Vegetation
	Building
	Low_Point
	Model_Key_Point
	Water
	Rail
	RoadSurface
	Overlap_Points
	Wire_Guard
	Wire_Conductor
	Transmission_Tower
	Wire_Structure_Connector
	Bridge_Deck
	High_Noise
	Overhead_Structure
	Ignored_Ground
	Snow
	Temporal_Exclusion
)

type Format0 struct {
	X              int32
	Y              int32
	Z              int32
	Intensity      uint16
	Pulse          uint8
	Classification uint8
	ScanAngleRank  int8
	UserData       uint8
	PointSourceID  uint16
}

type PDR0 struct {
	Format0
	ExtraBytes []byte
}

func (f0 *Format0) GetReturnNumber() uint8 {
	return f0.Pulse & PDR0_RETURN_NUMBER_MASK
}

func (f0 *Format0) GetNumberOfReturns() uint8 {
	return (f0.Pulse & PDR0_NUMBER_OF_RETURNS_MASK) >> 3
}

func (f0 *Format0) GetScanDirectionFlag() uint8 {
	return (f0.Pulse & PDR0_SCAN_DIRECTION_FLAG_MASK) >> 6
}

func (f0 *Format0) GetEdgeOfFlightLine() uint8 {
	return (f0.Pulse & PDR0_EDGE_OF_FLIGHT_LINE_MASK) >> 7
}

func (f0 *Format0) GetClassAttribute() ClassAttribute {
	return ClassAttribute(uint8(f0.Classification) & 0x1F)
}

// IsSynthetic if set, this point was created by a technique other than direct observation such as digitized from a photogrammetric
// stereo model or by traversing a waveform. Point attribute interpretation might differ from non-Synthetic points.
// Unused attributes must be set to the appropriate default value.
func (f0 *Format0) IsSynthetic() bool {
	return (uint8(f0.Classification) & 0x20) != 0
}

// IsKeyPoint if set, this point is considered to be a model keypoint and therefore generally should not be withheld in a
// thinning algorithm.
func (f0 *Format0) IsKeyPoint() bool {
	return (uint8(f0.Classification) & 0x40) != 0
}

// IsWithheld if set, this point should not be included in processing (synonymous with Deleted).
func (f0 *Format0) IsWithheld() bool {
	return (uint8(f0.Classification) & 0x80) != 0
}

type PDR0s []PDR0

func (p0 PDR0s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p0))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p0 {
		pdrSize := uint64(binary.Size(Format0{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p0[index].Format0); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p0[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p0[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p0 PDR0s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p0 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: 0,
			G: 0,
			B: 0,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______   __
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __| /_ |
// | |__  | |  | || |__) || \  / |   /  \   | |     | |
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |     | |
// | |    | |__| || | \ \ | |  | | / ____ \ | |     | |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|     |_|
//
//

type Format1 struct {
	Format0
	GPSTime float64
}

type PDR1 struct {
	Format1
	ExtraBytes []byte
}

type PDR1s []PDR1

func (p1 PDR1s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p1))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p1 {
		pdrSize := uint64(binary.Size(Format1{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p1[index].Format1); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p1[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p1[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p1 PDR1s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p1 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: 0,
			G: 0,
			B: 0,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______   ___
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __| |__ \
// | |__  | |  | || |__) || \  / |   /  \   | |       ) |
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |      / /
// | |    | |__| || | \ \ | |  | | / ____ \ | |     / /_
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|    |____|
//
//

type Format2 struct {
	Format0
	Red   uint16
	Green uint16
	Blue  uint16
}

type PDR2 struct {
	Format2
	ExtraBytes []byte
}

type PDR2s []PDR2

func (p2 PDR2s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p2))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p2 {
		pdrSize := uint64(binary.Size(Format2{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p2[index].Format2); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p2[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p2[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p2 PDR2s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p2 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: p.Red,
			G: p.Green,
			B: p.Blue,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______   ____
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __| |___ \
// | |__  | |  | || |__) || \  / |   /  \   | |      __) |
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |     |__ <
// | |    | |__| || | \ \ | |  | | / ____ \ | |     ___) |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|    |____/
//
//

type Format3 struct {
	Format1
	Red   uint16
	Green uint16
	Blue  uint16
}

type PDR3 struct {
	Format3
	ExtraBytes []byte
}

type PDR3s []PDR3

func (p3 PDR3s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p3))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p3 {
		pdrSize := uint64(binary.Size(Format3{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p3[index].Format3); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p3[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p3[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p3 PDR3s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p3 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: p.Red,
			G: p.Green,
			B: p.Blue,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______   _  _
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __| | || |
// | |__  | |  | || |__) || \  / |   /  \   | |    | || |_
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |    |__   _|
// | |    | |__| || | \ \ | |  | | / ____ \ | |       | |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|       |_|
//

type Format4 struct {
	Format1
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR4 struct {
	Format4
	ExtraBytes []byte
}

type PDR4s []PDR4

func (p4 PDR4s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p4))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p4 {
		pdrSize := uint64(binary.Size(Format4{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p4[index].Format4); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p4[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p4[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p4 PDR4s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p4 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: 0,
			G: 0,
			B: 0,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______   _____
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __| | ____|
// | |__  | |  | || |__) || \  / |   /  \   | |    | |__
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |    |___ \
// | |    | |__| || | \ \ | |  | | / ____ \ | |     ___) |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|    |____/
//
//

type Format5 struct {
	Format3
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR5 struct {
	Format5
	ExtraBytes []byte
}

type PDR5s []PDR5

func (p5 PDR5s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p5))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p5 {
		pdrSize := uint64(binary.Size(Format5{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p5[index].Format5); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p5[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p5[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p5 PDR5s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p5 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: p.Red,
			G: p.Green,
			B: p.Blue,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______     __
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __|   / /
// | |__  | |  | || |__) || \  / |   /  \   | |     / /_
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |    | '_ \
// | |    | |__| || | \ \ | |  | | / ____ \ | |    | (_) |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|     \___/
//
//

const (
	PDR6_RETURN_NUMBER_MASK        = 0x0F
	PDR6_NUMBER_OF_RETURNS_MASK    = 0xF0
	PDR6_CLASSIFICATION_FLAGS_MASK = 0x0F
	PDR6_SCAN_DIRECTION_FLAG_MASK  = 0x03
	PDR6_EDGE_OF_FLIGHT_LINE_MASK  = 0x0C
)

type Format6 struct {
	X              int32
	Y              int32
	Z              int32
	Intensity      uint16
	PulseReturns   uint8
	PulseFlags     uint8
	Classification uint8
	UserData       uint8
	ScanAngleRank  int16
	PointSourceID  uint16
	GPSTime        float64
}

type PDR6 struct {
	Format6
	ExtraBytes []byte
}

func (f6 *Format6) GetReturnNumber() uint8 {
	return f6.PulseReturns & PDR6_RETURN_NUMBER_MASK
}

func (f6 *Format6) GetNumberOfReturns() uint8 {
	return (f6.PulseReturns & PDR6_NUMBER_OF_RETURNS_MASK) >> 4
}

func (f6 *Format6) GetClassificationFlag() uint8 {
	return f6.PulseFlags & PDR6_CLASSIFICATION_FLAGS_MASK
}

func (f6 *Format6) GetScanDirectionFlag() uint8 {
	return (f6.PulseFlags & PDR6_SCAN_DIRECTION_FLAG_MASK) >> 4
}

func (f6 *Format6) GetEdgeOfFlightLine() uint8 {
	return (f6.PulseFlags & PDR6_EDGE_OF_FLIGHT_LINE_MASK) >> 2
}

func (f6 *Format6) GetClassAttribute() ClassAttribute {
	return ClassAttribute(uint8(f6.Classification) & 0x1F)
}

// IsSynthetic if set, this point was created by a technique other than direct observation such as digitized from a photogrammetric
// stereo model or by traversing a waveform. Point attribute interpretation might differ from non-Synthetic points.
// Unused attributes must be set to the appropriate default value.
func (f6 *Format6) IsSynthetic() bool {
	return (uint8(f6.Classification) & 0x20) != 0
}

// IsKeyPoint if set, this point is considered to be a model keypoint and therefore generally should not be withheld in a
// thinning algorithm.
func (f6 *Format6) IsKeyPoint() bool {
	return (uint8(f6.Classification) & 0x40) != 0
}

// IsWithheld if set, this point should not be included in processing (synonymous with Deleted).
func (f6 *Format6) IsWithheld() bool {
	return (uint8(f6.Classification) & 0x80) != 0
}

type PDR6s []PDR6

func (p6 PDR6s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p6))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p6 {
		pdrSize := uint64(binary.Size(Format6{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p6[index].Format6); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p6[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p6[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p6 PDR6s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p6 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: 0,
			G: 0,
			B: 0,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______   ______
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __| |____  |
// | |__  | |  | || |__) || \  / |   /  \   | |        / /
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |       / /
// | |    | |__| || | \ \ | |  | | / ____ \ | |      / /
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|     /_/
//
//

type Format7 struct {
	Format6
	Red   uint16
	Green uint16
	Blue  uint16
}

type PDR7 struct {
	Format7
	ExtraBytes []byte
}

type PDR7s []PDR7

func (p7 PDR7s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p7))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p7 {
		pdrSize := uint64(binary.Size(Format7{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p7[index].Format7); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p7[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p7[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p7 PDR7s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p7 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: p.Red,
			G: p.Green,
			B: p.Blue,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______    ___
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __|  / _ \
// | |__  | |  | || |__) || \  / |   /  \   | |    | (_) |
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |     > _ <
// | |    | |__| || | \ \ | |  | | / ____ \ | |    | (_) |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|     \___/
//
//

type Format8 struct {
	Format7
	NIR uint16
}

type PDR8 struct {
	Format8
	ExtraBytes []byte
}

type PDR8s []PDR8

func (p8 PDR8s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p8))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p8 {
		pdrSize := uint64(binary.Size(Format8{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p8[index].Format8); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p8[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p8[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p8 PDR8s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p8 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: p.Red,
			G: p.Green,
			B: p.Blue,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______    ___
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __|  / _ \
// | |__  | |  | || |__) || \  / |   /  \   | |    | (_) |
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |     \__, |
// | |    | |__| || | \ \ | |  | | / ____ \ | |       / /
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|      /_/
//
//

type Format9 struct {
	Format6
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR9 struct {
	Format9
	ExtraBytes []byte
}

type PDR9s []PDR9

func (p9 PDR9s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p9))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p9 {
		pdrSize := uint64(binary.Size(Format9{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p9[index].Format9); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p9[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p9[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p9 PDR9s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p9 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: 0,
			G: 0,
			B: 0,
		}
		output = append(output, csvRow)
	}
	return
}

//  ______  ____   _____   __  __         _______   __   ___
// |  ____|/ __ \ |  __ \ |  \/  |    /\ |__   __| /_ | / _ \
// | |__  | |  | || |__) || \  / |   /  \   | |     | || | | |
// |  __| | |  | ||  _  / | |\/| |  / /\ \  | |     | || | | |
// | |    | |__| || | \ \ | |  | | / ____ \ | |     | || |_| |
// |_|     \____/ |_|  \_\|_|  |_|/_/    \_\|_|     |_| \___/
//
//

type Format10 struct {
	Format8
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR10 struct {
	Format10
	ExtraBytes []byte
}

type PDR10s []PDR10

func (p10 PDR10s) read(file *os.File, offsetIn int64, dataLength uint64) (err error) {
	bytesToRead := make([]byte, uint64(len(p10))*dataLength)
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	for index := range p10 {
		pdrSize := uint64(binary.Size(Format10{}))
		offset := uint64(index) * dataLength
		if err = binary.Read(bytes.NewReader(bytesToRead[offset:offset+pdrSize]), binary.LittleEndian, &p10[index].Format10); err != nil {
			return
		}
		if dataLength-pdrSize != 0 {
			p10[index].ExtraBytes = make([]byte, dataLength-pdrSize)
			if err = binary.Read(bytes.NewReader(bytesToRead[offset+pdrSize:offset+dataLength]), binary.LittleEndian, &p10[index].ExtraBytes); err != nil {
				return
			}
		}
	}
	return
}

func (p10 PDR10s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p10 {
		csvRow := &XYZRGB{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
			R: p.Red,
			G: p.Green,
			B: p.Blue,
		}
		output = append(output, csvRow)
	}
	return
}
