package IPMTSCD

import "github.com/TOPAS-2545/ber"

// IpmstscdData is the top-level type (replaces IPMSTSCDData)
type IpmstscdData struct {
	DetectorControllerInfo *DetectorControllerInfo `asn1:"optional,tag:0"`
	DetectorInfo           []DetectorInfoSeq       `asn1:"optional,tag:1"`
}

// DetectorControllerInfo (replaces DetectorControllerInformation)
type DetectorControllerInfo struct {
	DetectorControllerIndex        int64                    `asn1:"tag:0"`
	DetectorControllerTimeLocation *GeneralTimeLocationCore `asn1:"optional,tag:1"`
}

// GeneralTimeLocationCore represents time and location information.
// Optional integer fields are pointer types so that "present with value 0"
// is distinguishable from "absent" across the BER encode/decode round-trip.
type GeneralTimeLocationCore struct {
	CurrentTime       Time   `asn1:"tag:0"`
	LocationLongitude *int64 `asn1:"optional,tag:1"`
	LocationLatitude  *int64 `asn1:"optional,tag:2"`
	LocationElevation *int64 `asn1:"optional,tag:3"`
}

// DetectorInfoSeq (replaces IpmstscdDetDataSeq)
type DetectorInfoSeq struct {
	PhysicalDetectorIndex int64                    `asn1:"tag:0"`
	InfoType              ber.Enumerated           `asn1:"tag:1"`
	DetectionInfo         ber.RawValue             `asn1:"explicit,tag:2"`
	DetectorTimeLocation  *GeneralTimeLocationCore `asn1:"optional,tag:3"`
}

// OccupancyBasedDetectionData (replaces IpmstscdOccTypeDetectorInformation)
// CHOICE tag [1]
type OccupancyBasedDetectionData struct {
	OccupancyDataDuration          *int64                         `asn1:"optional,tag:0"`
	OccupancyState                 bool                           `asn1:"tag:1"`
	OccupancyStateDuration         int64                          `asn1:"tag:2"`
	OccupancyPreviousStateDuration int64                          `asn1:"tag:3"`
	OccupancyRate                  float64                        `asn1:"tag:4"`
	Speed                          *float64                       `asn1:"optional,tag:5"`
	Volume                         int64                          `asn1:"tag:6"`
	OccupancyNonOccupancyHistory   []OccupancyNonOccupancyHistory `asn1:"optional,tag:7"`
	ErrorState                     *ber.Enumerated                `asn1:"optional,tag:8"`
	UserData                       []byte                         `asn1:"optional,tag:9"`
	TargetType                     *int64                         `asn1:"optional,tag:10"`
	DirectionDiscrimination        *bool                          `asn1:"optional,tag:11"`
}

// ImageProcessingBasedDetectionData (replaces IpmstscdImageTypeDetectorInformation)
// CHOICE tag [2]
type ImageProcessingBasedDetectionData struct {
	DataDuration                 *int64                        `asn1:"optional,tag:0"`
	VolumeQuantity               int64                         `asn1:"tag:1"`
	QueueLength                  *int64                        `asn1:"optional,tag:2"`
	OccupancyRate                *float64                      `asn1:"optional,tag:3"`
	SpeedQuantity                *float64                      `asn1:"optional,tag:4"`
	OccupancyNonOccupancyHistory *OccupancyNonOccupancyHistory `asn1:"optional,tag:5"`
	ErrorState                   *ber.Enumerated               `asn1:"optional,tag:6"`
	UserData                     []byte                        `asn1:"optional,tag:7"`
}

// IdentificationBasedDetectionData (replaces IpmstscdIDTypeDetectorInformation)
// CHOICE tag [3]
type IdentificationBasedDetectionData struct {
	SequenceNumber      int64           `asn1:"tag:0"`
	DeviceType          *ber.Enumerated `asn1:"optional,tag:1"`
	VehicleId           []byte          `asn1:"tag:2"`
	VehicleType         *int64          `asn1:"optional,tag:3"`
	VehicleUse          *int64          `asn1:"optional,tag:4"`
	DetectionLane       *int64          `asn1:"optional,tag:5"`
	DetectionLaneMedian *int64          `asn1:"optional,tag:6"`
	DetectionSpeed      *float64        `asn1:"optional,tag:7"`
	Occupancy           *int64          `asn1:"optional,tag:8"`
	ErrorState          *ber.Enumerated `asn1:"optional,tag:9"`
	TagInfo             []byte          `asn1:"optional,tag:10"`
	UserData            []byte          `asn1:"optional,tag:11"`
}

// OccupancyBasedDetectionDataType2 - Type2 variant
// CHOICE tag [4]
type OccupancyBasedDetectionDataType2 struct {
	AccumulativeDetection    []AccumulativeDetectionSeq    `asn1:"tag:0"`
	PassingVehicleTimeSeries []PassingVehicleTimeSeriesSeq `asn1:"tag:1"`
	VehicleSpeed             []VehicleSpeedSeq             `asn1:"tag:2"`
	SpecificVehicleDetection []byte                        `asn1:"tag:3"`
	DetectorStatus           *ber.BitString                `asn1:"optional,tag:4"`
}

// ImageProcessingBasedDetectionDataType2 - Type2 variant
// CHOICE tag [5]
type ImageProcessingBasedDetectionDataType2 struct {
	AccumulativeDetection    []AccumulativeDetectionSeq    `asn1:"tag:0"`
	PassingVehicleTimeSeries []PassingVehicleTimeSeriesSeq `asn1:"tag:1"`
	VehicleSpeed             []VehicleSpeedSeq             `asn1:"tag:2"`
	SpecificVehicleDetection []byte                        `asn1:"tag:3"`
	CongestionInfo           CongestionInfo                `asn1:"tag:4"`
	DirectionalTrafficVolume []DirectionalTrafficVolumeSeq `asn1:"tag:5"`
	DetectorStatus           *ber.BitString                `asn1:"optional,tag:6"`
}

// CongestionInfo holds queue length and start position info
type CongestionInfo struct {
	CongestionLength1     int64 `asn1:"tag:0"`
	VehicleStartPosition1 int64 `asn1:"tag:1"`
	CongestionLength2     int64 `asn1:"tag:2"`
	VehicleStartPosition2 int64 `asn1:"tag:3"`
}

// IdentificationBasedDetectionDataType2 - Type2 variant
// CHOICE tag [6]
type IdentificationBasedDetectionDataType2 struct {
	AccumulativeDetection    []AccumulativeDetectionSeq    `asn1:"tag:0"`
	PassingVehicleTimeSeries []PassingVehicleTimeSeriesSeq `asn1:"tag:1"`
	VehicleSpeed             []VehicleSpeedSeq             `asn1:"tag:2"`
	SpecificVehicleDetection []byte                        `asn1:"tag:3"`
	VehicleIdentification    []VehicleIdentificationSeq    `asn1:"tag:4"`
	DetectorStatus           *ber.BitString                `asn1:"optional,tag:5"`
}

// OccupancyNonOccupancyHistory represents occupancy/non-occupancy pair history
type OccupancyNonOccupancyHistory struct {
	OccupancyTimes    int64 `asn1:"tag:0"`
	NonOccupancyTimes int64 `asn1:"tag:1"`
}

// Time represents the ASN.1 Time type
type Time struct {
	Year            *int64       `asn1:"optional,tag:0"`
	Month           *int64       `asn1:"optional,tag:1"`
	Day             *int64       `asn1:"optional,tag:2"`
	Hour            *int64       `asn1:"optional,tag:3"`
	Minute          *int64       `asn1:"optional,tag:4"`
	Second          *int64       `asn1:"optional,tag:5"`
	SecondFractions ber.RawValue `asn1:"tag:6"`
	TimeZone        *TimeZone    `asn1:"optional,tag:7"`
}

// TimeZone represents the timezone sub-sequence within Time
type TimeZone struct {
	Hour   *int64 `asn1:"optional,tag:0"`
	Minute *int64 `asn1:"optional,tag:1"`
}

// AccumulativeDetectionSeq represents accumulative detection data per detector
type AccumulativeDetectionSeq struct {
	DetectorNumber     int64           `asn1:"tag:0"`
	DetectionStatus    *ber.Enumerated `asn1:"optional,tag:1"`
	Density            int64           `asn1:"tag:2"`
	Occupancy          int64           `asn1:"tag:3"`
	DetectorPulseError int64           `asn1:"tag:4"`
}

// PassingVehicleTimeSeriesSeq represents passing vehicle time series data per detector
type PassingVehicleTimeSeriesSeq struct {
	DetectorNumber        int64           `asn1:"tag:0"`
	DetectionStatus       *ber.Enumerated `asn1:"optional,tag:1"`
	VehiclePassageHistory []byte          `asn1:"tag:2"`
}

// VehicleSpeedSeq represents vehicle speed data per detector
type VehicleSpeedSeq struct {
	DetectorNumber int64          `asn1:"tag:0"`
	VehicleType    ber.Enumerated `asn1:"tag:1"`
	Velocity       int64          `asn1:"tag:2"`
}

// VehicleIdentificationSeq represents vehicle identification data
type VehicleIdentificationSeq struct {
	VehicleID []byte `asn1:"tag:0"`
	Data      []byte `asn1:"optional,tag:1"`
}

// DirectionalTrafficVolumeSeq represents directional traffic volume data
type DirectionalTrafficVolumeSeq struct {
	DirectionNumber  int64           `asn1:"tag:0"`
	DetectionStatus  *ber.Enumerated `asn1:"optional,tag:1"`
	DirectionDensity int64           `asn1:"tag:2"`
}
