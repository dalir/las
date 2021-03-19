package main

import "fmt"

type LasSepcVersion string

const (
	V1_4 LasSepcVersion = "1.4"
	V1_3                = "1.3"
	V1_2                = "1.2"
	V1_1                = "1.1"
)

var AllLasVersions = []LasSepcVersion{V1_4, V1_3, V1_2, V1_1}

//  _____  _    _ ____  _      _____ _____   _    _ ______          _____  ______ _____    ____  _      ____   _____ _  __
// |  __ \| |  | |  _ \| |    |_   _/ ____| | |  | |  ____|   /\   |  __ \|  ____|  __ \  |  _ \| |    / __ \ / ____| |/ /
// | |__) | |  | | |_) | |      | || |      | |__| | |__     /  \  | |  | | |__  | |__) | | |_) | |   | |  | | |    | ' /
// |  ___/| |  | |  _ <| |      | || |      |  __  |  __|   / /\ \ | |  | |  __| |  _  /  |  _ <| |   | |  | | |    |  <
// | |    | |__| | |_) | |____ _| || |____  | |  | | |____ / ____ \| |__| | |____| | \ \  | |_) | |___| |__| | |____| . \
// |_|     \____/|____/|______|_____\_____| |_|  |_|______/_/    \_\_____/|______|_|  \_\ |____/|______\____/ \_____|_|\_\
//
//

type PublicHeaderBlock struct {
	FileSignature                            [4]byte
	FileSourceID                             uint16
	GlobalEncoding                           uint16
	GUID1                                    uint32
	GUID2                                    uint16
	GUID3                                    uint16
	GUID4                                    [8]byte
	VersionMajor                             byte
	VersionMinor                             byte
	SystemID                                 [32]byte
	GeneratingSoftware                       [32]byte
	FileCreationDayOfYear                    uint16
	FileCreationYear                         uint16
	HeaderSize                               uint16
	OffsetToPointData                        uint32
	NumberOfVLRs                             uint32
	PointDataRecordFormat                    byte
	PointDataRecordLength                    uint16
	LegacyNumberOfPointRecords               uint32
	LegacyNumberOfPointByReturn              [5]uint32
	XScaleFactor                             float64
	YScaleFactor                             float64
	ZScaleFactor                             float64
	XOffset                                  float64
	YOffset                                  float64
	ZOffset                                  float64
	MaxX                                     float64
	MinX                                     float64
	MaxY                                     float64
	MinY                                     float64
	MaxZ                                     float64
	MinZ                                     float64
	StartOfWaveformDataPacketRecord          uint64
	StartOfFirstExtendedVariableLengthRecord uint64
	NumberOfExtendedVariableLengthRecords    uint32
	NumberOfPointRecords                     uint64
	NumberOfPointsByReturn                   [15]uint64
}

func (phb *PublicHeaderBlock) GetVersion() (version LasSepcVersion) {
	version = LasSepcVersion(fmt.Sprintf("%d.%d", phb.VersionMajor, phb.VersionMinor))
	return
}
