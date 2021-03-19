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
type PDR interface {
	read(file *os.File, offsetIn int64) (offsetOut int64, err error)
	GetXYZRGB() (output XYZRGB)
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

func (p0 *PDR0) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR0{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p0); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR0{}))
	return
}

func (p0 *PDR0) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p0.X),
		Y: float64(p0.Y),
		Z: float64(p0.Z),
		R: 0,
		G: 0,
		B: 0,
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

func (p1 *PDR1) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR1{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p1); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR1{}))
	return
}

func (p1 *PDR1) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p1.Format0.X),
		Y: float64(p1.Format0.Y),
		Z: float64(p1.Format0.Z),
		R: 0,
		G: 0,
		B: 0,
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

func (p2 *PDR2) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR2{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p2); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR2{}))
	return
}

func (p2 *PDR2) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p2.Format0.X),
		Y: float64(p2.Format0.Y),
		Z: float64(p2.Format0.Z),
		R: p2.Red,
		G: p2.Green,
		B: p2.Blue,
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

func (p3 *PDR3) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR3{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p3); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR3{}))
	return
}

func (p3 *PDR3) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p3.Format1.Format0.X),
		Y: float64(p3.Format1.Format0.Y),
		Z: float64(p3.Format1.Format0.Z),
		R: p3.Red,
		G: p3.Green,
		B: p3.Blue,
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

func (p4 *PDR4) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR4{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p4); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR4{}))
	return
}

func (p4 *PDR4) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p4.Format1.Format0.X),
		Y: float64(p4.Format1.Format0.Y),
		Z: float64(p4.Format1.Format0.Z),
		R: 0,
		G: 0,
		B: 0,
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

func (p5 *PDR5) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR5{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p5); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR5{}))
	return
}

func (p5 *PDR5) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p5.Format3.Format1.Format0.X),
		Y: float64(p5.Format3.Format1.Format0.Y),
		Z: float64(p5.Format3.Format1.Format0.Z),
		R: p5.Format3.Red,
		G: p5.Format3.Green,
		B: p5.Format3.Blue,
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

func (p *PDR6) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR6{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR6{}))
	return
}

func (p6 *PDR6) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p6.X),
		Y: float64(p6.Y),
		Z: float64(p6.Z),
		R: 0,
		G: 0,
		B: 0,
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

func (p7 *PDR7) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR7{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p7); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR7{}))
	return
}

func (p7 *PDR7) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p7.Format6.X),
		Y: float64(p7.Format6.Y),
		Z: float64(p7.Format6.Z),
		R: p7.Red,
		G: p7.Green,
		B: p7.Blue,
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

func (p8 *PDR8) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR8{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p8); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR8{}))
	return
}

func (p8 *PDR8) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p8.Format7.Format6.X),
		Y: float64(p8.Format7.Format6.Y),
		Z: float64(p8.Format7.Format6.Z),
		R: p8.Format7.Red,
		G: p8.Format7.Green,
		B: p8.Format7.Blue,
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

func (p9 *PDR9) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR9{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p9); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR9{}))
	return
}

func (p9 *PDR9) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p9.Format6.X),
		Y: float64(p9.Format6.Y),
		Z: float64(p9.Format6.Z),
		R: 0,
		G: 0,
		B: 0,
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

func (p10 *PDR10) read(file *os.File, offsetIn int64) (offsetOut int64, err error) {
	pdrInBytes := make([]byte, binary.Size(PDR10{}))
	_, err = file.ReadAt(pdrInBytes, offsetIn)
	if err != nil {
		return
	}
	if err = binary.Read(bytes.NewReader(pdrInBytes), binary.LittleEndian, p10); err != nil {
		return
	}
	offsetOut = offsetIn + int64(binary.Size(PDR10{}))
	return
}

func (p10 *PDR10) GetXYZRGB() (output XYZRGB) {
	output = XYZRGB{
		X: float64(p10.Format8.Format7.Format6.X),
		Y: float64(p10.Format8.Format7.Format6.Y),
		Z: float64(p10.Format8.Format7.Format6.Z),
		R: p10.Format8.Format7.Red,
		G: p10.Format8.Format7.Green,
		B: p10.Format8.Format7.Blue,
	}
	return
}
