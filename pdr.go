package main

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
	read(file *os.File, offsetIn int64) (err error)
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

type PDR0 struct {
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

func (p0 *PDR0) GetReturnNumber() uint8 {
	return p0.Pulse & PDR0_RETURN_NUMBER_MASK
}

func (p0 *PDR0) GetNumberOfReturns() uint8 {
	return (p0.Pulse & PDR0_NUMBER_OF_RETURNS_MASK) >> 3
}

func (p0 *PDR0) GetScanDirectionFlag() uint8 {
	return (p0.Pulse & PDR0_SCAN_DIRECTION_FLAG_MASK) >> 6
}

func (p0 *PDR0) GetEdgeOfFlightLine() uint8 {
	return (p0.Pulse & PDR0_EDGE_OF_FLIGHT_LINE_MASK) >> 7
}

type PDR0s []PDR0

func (p0 PDR0s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p0)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p0); err != nil {
		return
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

type PDR1 struct {
	Format0 PDR0
	GPSTime float64
}

type PDR1s []PDR1

func (p1 PDR1s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p1)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p1); err != nil {
		return
	}
	return
}

func (p1 PDR1s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p1 {
		csvRow := &XYZRGB{
			X: float64(p.Format0.X),
			Y: float64(p.Format0.Y),
			Z: float64(p.Format0.Z),
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

type PDR2 struct {
	Format0 PDR0
	Red     uint16
	Green   uint16
	Blue    uint16
}

type PDR2s []PDR2

func (p2 PDR2s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p2)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p2); err != nil {
		return
	}
	return
}

func (p2 PDR2s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p2 {
		csvRow := &XYZRGB{
			X: float64(p.Format0.X),
			Y: float64(p.Format0.Y),
			Z: float64(p.Format0.Z),
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

type PDR3 struct {
	Format1 PDR1
	Red     uint16
	Green   uint16
	Blue    uint16
}

type PDR3s []PDR3

func (p3 PDR3s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p3)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p3); err != nil {
		return
	}
	return
}

func (p3 PDR3s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p3 {
		csvRow := &XYZRGB{
			X: float64(p.Format1.Format0.X),
			Y: float64(p.Format1.Format0.Y),
			Z: float64(p.Format1.Format0.Z),
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

type PDR4 struct {
	Format1                     PDR1
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR4s []PDR4

func (p4 PDR4s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p4)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p4); err != nil {
		return
	}
	return
}

func (p4 PDR4s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p4 {
		csvRow := &XYZRGB{
			X: float64(p.Format1.Format0.X),
			Y: float64(p.Format1.Format0.Y),
			Z: float64(p.Format1.Format0.Z),
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

type PDR5 struct {
	Format3                     PDR3
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR5s []PDR5

func (p5 PDR5s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p5)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p5); err != nil {
		return
	}
	return
}

func (p5 PDR5s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p5 {
		csvRow := &XYZRGB{
			X: float64(p.Format3.Format1.Format0.X),
			Y: float64(p.Format3.Format1.Format0.Y),
			Z: float64(p.Format3.Format1.Format0.Z),
			R: p.Format3.Red,
			G: p.Format3.Green,
			B: p.Format3.Blue,
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

type PDR6 struct {
	X              int32
	Y              int32
	Z              int32
	Intensity      uint16
	PulseReturns   uint8
	PulseFlags     uint8
	Classification uint8
	ScanAngleRank  int8
	UserData       uint8
	PointSourceID  uint16
	GPSTime        float64
}

func (p6 *PDR6) GetReturnNumber() uint8 {
	return p6.PulseReturns & PDR6_RETURN_NUMBER_MASK
}

func (p6 *PDR6) GetNumberOfReturns() uint8 {
	return (p6.PulseReturns & PDR6_NUMBER_OF_RETURNS_MASK) >> 4
}

func (p6 *PDR6) GetClassificationFlag() uint8 {
	return p6.PulseFlags & PDR6_CLASSIFICATION_FLAGS_MASK
}

func (p6 *PDR6) GetScanDirectionFlag() uint8 {
	return (p6.PulseFlags & PDR6_SCAN_DIRECTION_FLAG_MASK) >> 4
}

func (p6 *PDR6) GetEdgeOfFlightLine() uint8 {
	return (p6.PulseFlags & PDR6_EDGE_OF_FLIGHT_LINE_MASK) >> 2
}

type PDR6s []PDR6

func (p6 PDR6s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p6)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p6); err != nil {
		return
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

type PDR7 struct {
	Format6 PDR6
	Red     uint16
	Green   uint16
	Blue    uint16
}

type PDR7s []PDR7

func (p7 PDR7s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p7)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p7); err != nil {
		return
	}
	return
}

func (p7 PDR7s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p7 {
		csvRow := &XYZRGB{
			X: float64(p.Format6.X),
			Y: float64(p.Format6.Y),
			Z: float64(p.Format6.Z),
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

type PDR8 struct {
	Format7 PDR7
	NIR     uint16
}

type PDR8s []PDR8

func (p8 PDR8s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p8)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p8); err != nil {
		return
	}
	return
}

func (p8 PDR8s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p8 {
		csvRow := &XYZRGB{
			X: float64(p.Format7.Format6.X),
			Y: float64(p.Format7.Format6.Y),
			Z: float64(p.Format7.Format6.Z),
			R: p.Format7.Red,
			G: p.Format7.Green,
			B: p.Format7.Blue,
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

type PDR9 struct {
	Format6                     PDR6
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR9s []PDR9

func (p9 PDR9s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p9)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p9); err != nil {
		return
	}
	return
}

func (p9 PDR9s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p9 {
		csvRow := &XYZRGB{
			X: float64(p.Format6.X),
			Y: float64(p.Format6.Y),
			Z: float64(p.Format6.Z),
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

type PDR10 struct {
	Format8                     PDR8
	WavePacketDescriptorIndex   uint8
	ByteOffsetToWaveformData    uint64
	WaveformPacketSizeInBytes   uint32
	ReturnPointWaveformLocation float32
	ParametricDx                float32
	ParametricDy                float32
	ParametricDz                float32
}

type PDR10s []PDR10

func (p10 PDR10s) read(file *os.File, offsetIn int64) (err error) {
	bytesToRead := make([]byte, uint64(binary.Size(p10)))
	_, err = file.ReadAt(bytesToRead, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(bytesToRead), binary.LittleEndian, p10); err != nil {
		return
	}
	return
}

func (p10 PDR10s) GetCSVList() (output []*XYZRGB) {
	for _, p := range p10 {
		csvRow := &XYZRGB{
			X: float64(p.Format8.Format7.Format6.X),
			Y: float64(p.Format8.Format7.Format6.Y),
			Z: float64(p.Format8.Format7.Format6.Z),
			R: p.Format8.Format7.Red,
			G: p.Format8.Format7.Green,
			B: p.Format8.Format7.Blue,
		}
		output = append(output, csvRow)
	}
	return
}
