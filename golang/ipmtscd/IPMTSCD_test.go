package IPMTSCD

// Feature: asn1-spec-update, Property 1: Go BER CHOICE round-trip preservation
// **Validates: Requirements 3.8, 4.3**
//
// For any valid DetectionInfo CHOICE variant (tags 1-6), marshalling to BER
// and unmarshalling back should produce an equal struct.

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/TOPAS-2545/ber"
)

const propertyTestIterations = 100

// --- Helpers for pointer fields ---

func ptrInt64(v int64) *int64                         { return &v }
func ptrFloat64(v float64) *float64                   { return &v }
func ptrEnum(v ber.Enumerated) *ber.Enumerated        { return &v }

// --- Generators for each CHOICE variant ---

func genOccupancyBasedDetectionData(rng *rand.Rand) OccupancyBasedDetectionData {
	d := OccupancyBasedDetectionData{
		OccupancyState:                 rng.Intn(2) == 1,
		OccupancyStateDuration:         int64(rng.Intn(65536)),
		OccupancyPreviousStateDuration: int64(rng.Intn(65536)),
		OccupancyRate:                  float64(rng.Intn(1000)) / 10.0,
		Volume:                         int64(rng.Intn(10000)),
	}
	// Optionally set OccupancyDataDuration
	if rng.Intn(2) == 1 {
		d.OccupancyDataDuration = ptrInt64(int64(rng.Intn(3600)))
	}
	// Optionally set Speed
	if rng.Intn(2) == 1 {
		d.Speed = ptrFloat64(float64(rng.Intn(200)) + 0.5)
	}
	return d
}

func genImageProcessingBasedDetectionData(rng *rand.Rand) ImageProcessingBasedDetectionData {
	d := ImageProcessingBasedDetectionData{
		VolumeQuantity: int64(rng.Intn(10000)),
	}
	if rng.Intn(2) == 1 {
		d.DataDuration = ptrInt64(int64(rng.Intn(3600)))
	}
	if rng.Intn(2) == 1 {
		d.QueueLength = ptrInt64(int64(rng.Intn(500)))
	}
	if rng.Intn(2) == 1 {
		d.OccupancyRate = ptrFloat64(float64(rng.Intn(1000)) / 10.0)
	}
	if rng.Intn(2) == 1 {
		d.SpeedQuantity = ptrFloat64(float64(rng.Intn(200)) + 0.5)
	}
	return d
}

func genIdentificationBasedDetectionData(rng *rand.Rand) IdentificationBasedDetectionData {
	vid := make([]byte, 1+rng.Intn(8))
	for i := range vid {
		vid[i] = byte(rng.Intn(256))
	}
	d := IdentificationBasedDetectionData{
		SequenceNumber: int64(rng.Intn(256)),
		VehicleId:      vid,
	}
	if rng.Intn(2) == 1 {
		d.DeviceType = ptrEnum(ber.Enumerated(rng.Intn(5)))
	}
	if rng.Intn(2) == 1 {
		d.VehicleType = ptrInt64(int64(rng.Intn(20)))
	}
	if rng.Intn(2) == 1 {
		d.DetectionLane = ptrInt64(int64(1 + rng.Intn(8)))
	}
	if rng.Intn(2) == 1 {
		d.DetectionSpeed = ptrFloat64(float64(rng.Intn(200)) + 0.5)
	}
	return d
}

func genAccumulativeDetectionSeq(rng *rand.Rand) AccumulativeDetectionSeq {
	return AccumulativeDetectionSeq{
		DetectorNumber:     int64(rng.Intn(256)),
		Density:            int64(rng.Intn(100)),
		Occupancy:          int64(rng.Intn(100)),
		DetectorPulseError: int64(rng.Intn(100)),
	}
}

func genPassingVehicleTimeSeriesSeq(rng *rand.Rand) PassingVehicleTimeSeriesSeq {
	hist := make([]byte, 1+rng.Intn(4))
	for i := range hist {
		hist[i] = byte(rng.Intn(256))
	}
	return PassingVehicleTimeSeriesSeq{
		DetectorNumber:        int64(rng.Intn(256)),
		VehiclePassageHistory: hist,
	}
}

func genVehicleSpeedSeq(rng *rand.Rand) VehicleSpeedSeq {
	return VehicleSpeedSeq{
		DetectorNumber: int64(rng.Intn(256)),
		VehicleType:    ber.Enumerated(rng.Intn(5)),
		Velocity:       int64(rng.Intn(300)),
	}
}

func genOccupancyBasedDetectionDataType2(rng *rand.Rand) OccupancyBasedDetectionDataType2 {
	numAcc := 1 + rng.Intn(3)
	acc := make([]AccumulativeDetectionSeq, numAcc)
	for i := range acc {
		acc[i] = genAccumulativeDetectionSeq(rng)
	}
	numPV := 1 + rng.Intn(3)
	pv := make([]PassingVehicleTimeSeriesSeq, numPV)
	for i := range pv {
		pv[i] = genPassingVehicleTimeSeriesSeq(rng)
	}
	numVS := 1 + rng.Intn(3)
	vs := make([]VehicleSpeedSeq, numVS)
	for i := range vs {
		vs[i] = genVehicleSpeedSeq(rng)
	}
	svd := make([]byte, 1+rng.Intn(4))
	for i := range svd {
		svd[i] = byte(rng.Intn(256))
	}
	d := OccupancyBasedDetectionDataType2{
		AccumulativeDetection:    acc,
		PassingVehicleTimeSeries: pv,
		VehicleSpeed:             vs,
		SpecificVehicleDetection: svd,
	}
	if rng.Intn(2) == 1 {
		d.DetectorStatus = &ber.BitString{
			Bytes:     []byte{byte(rng.Intn(4)) << 6}, // 2 named bits
			BitLength: 2,
		}
	}
	return d
}

func genImageProcessingBasedDetectionDataType2(rng *rand.Rand) ImageProcessingBasedDetectionDataType2 {
	numAcc := 1 + rng.Intn(3)
	acc := make([]AccumulativeDetectionSeq, numAcc)
	for i := range acc {
		acc[i] = genAccumulativeDetectionSeq(rng)
	}
	numPV := 1 + rng.Intn(3)
	pv := make([]PassingVehicleTimeSeriesSeq, numPV)
	for i := range pv {
		pv[i] = genPassingVehicleTimeSeriesSeq(rng)
	}
	numVS := 1 + rng.Intn(3)
	vs := make([]VehicleSpeedSeq, numVS)
	for i := range vs {
		vs[i] = genVehicleSpeedSeq(rng)
	}
	svd := make([]byte, 1+rng.Intn(4))
	for i := range svd {
		svd[i] = byte(rng.Intn(256))
	}
	numDTV := 1 + rng.Intn(3)
	dtv := make([]DirectionalTrafficVolumeSeq, numDTV)
	for i := range dtv {
		dtv[i] = DirectionalTrafficVolumeSeq{
			DirectionNumber:  int64(rng.Intn(256)),
			DirectionDensity: int64(rng.Intn(100)),
		}
	}
	d := ImageProcessingBasedDetectionDataType2{
		AccumulativeDetection:    acc,
		PassingVehicleTimeSeries: pv,
		VehicleSpeed:             vs,
		SpecificVehicleDetection: svd,
		CongestionInfo: CongestionInfo{
			CongestionLength1:     int64(rng.Intn(1000)),
			VehicleStartPosition1: int64(rng.Intn(1000)),
			CongestionLength2:     int64(rng.Intn(1000)),
			VehicleStartPosition2: int64(rng.Intn(1000)),
		},
		DirectionalTrafficVolume: dtv,
	}
	if rng.Intn(2) == 1 {
		d.DetectorStatus = &ber.BitString{
			Bytes:     []byte{byte(rng.Intn(4)) << 6},
			BitLength: 2,
		}
	}
	return d
}

func genIdentificationBasedDetectionDataType2(rng *rand.Rand) IdentificationBasedDetectionDataType2 {
	numAcc := 1 + rng.Intn(3)
	acc := make([]AccumulativeDetectionSeq, numAcc)
	for i := range acc {
		acc[i] = genAccumulativeDetectionSeq(rng)
	}
	numPV := 1 + rng.Intn(3)
	pv := make([]PassingVehicleTimeSeriesSeq, numPV)
	for i := range pv {
		pv[i] = genPassingVehicleTimeSeriesSeq(rng)
	}
	numVS := 1 + rng.Intn(3)
	vs := make([]VehicleSpeedSeq, numVS)
	for i := range vs {
		vs[i] = genVehicleSpeedSeq(rng)
	}
	svd := make([]byte, 1+rng.Intn(4))
	for i := range svd {
		svd[i] = byte(rng.Intn(256))
	}
	numVI := 1 + rng.Intn(3)
	vi := make([]VehicleIdentificationSeq, numVI)
	for i := range vi {
		vid := make([]byte, 1+rng.Intn(8))
		for j := range vid {
			vid[j] = byte(rng.Intn(256))
		}
		vi[i] = VehicleIdentificationSeq{VehicleID: vid}
	}
	d := IdentificationBasedDetectionDataType2{
		AccumulativeDetection:    acc,
		PassingVehicleTimeSeries: pv,
		VehicleSpeed:             vs,
		SpecificVehicleDetection: svd,
		VehicleIdentification:    vi,
	}
	if rng.Intn(2) == 1 {
		d.DetectorStatus = &ber.BitString{
			Bytes:     []byte{byte(rng.Intn(4)) << 6},
			BitLength: 2,
		}
	}
	return d
}

// --- Round-trip helper ---

// marshalChoiceVariant marshals a CHOICE variant with the given context-specific tag.
func marshalChoiceVariant(val interface{}, tag int) ([]byte, error) {
	// First marshal the inner struct to get its BER encoding
	innerBytes, err := ber.Marshal(val)
	if err != nil {
		return nil, err
	}
	// Wrap in a RawValue with context-specific class and the CHOICE tag
	raw := ber.RawValue{
		Class:      ber.ClassContextSpecific,
		Tag:        tag,
		IsCompound: true,
		Bytes:      innerBytes,
	}
	return ber.Marshal(raw)
}

// unmarshalChoiceVariant unmarshals BER bytes (a context-specific tagged value)
// back into the target struct.
func unmarshalChoiceVariant(data []byte, target interface{}) error {
	// First unmarshal to get the RawValue (which contains the inner bytes)
	var raw ber.RawValue
	_, err := ber.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	// Then unmarshal the inner bytes into the target struct
	_, err = ber.Unmarshal(raw.Bytes, target)
	return err
}

// --- Property test ---

func TestPropertyBERChoiceRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < propertyTestIterations; i++ {
		// Pick a random CHOICE variant (tags 1-6)
		choiceTag := 1 + rng.Intn(6)

		switch choiceTag {
		case 1:
			original := genOccupancyBasedDetectionData(rng)
			encoded, err := marshalChoiceVariant(original, 1)
			if err != nil {
				t.Fatalf("iteration %d: marshal OccupancyBasedDetectionData failed: %v", i, err)
			}
			var decoded OccupancyBasedDetectionData
			err = unmarshalChoiceVariant(encoded, &decoded)
			if err != nil {
				t.Fatalf("iteration %d: unmarshal OccupancyBasedDetectionData failed: %v", i, err)
			}
			if !reflect.DeepEqual(original, decoded) {
				t.Fatalf("iteration %d: OccupancyBasedDetectionData round-trip mismatch\noriginal: %+v\ndecoded:  %+v", i, original, decoded)
			}

		case 2:
			original := genImageProcessingBasedDetectionData(rng)
			encoded, err := marshalChoiceVariant(original, 2)
			if err != nil {
				t.Fatalf("iteration %d: marshal ImageProcessingBasedDetectionData failed: %v", i, err)
			}
			var decoded ImageProcessingBasedDetectionData
			err = unmarshalChoiceVariant(encoded, &decoded)
			if err != nil {
				t.Fatalf("iteration %d: unmarshal ImageProcessingBasedDetectionData failed: %v", i, err)
			}
			if !reflect.DeepEqual(original, decoded) {
				t.Fatalf("iteration %d: ImageProcessingBasedDetectionData round-trip mismatch\noriginal: %+v\ndecoded:  %+v", i, original, decoded)
			}

		case 3:
			original := genIdentificationBasedDetectionData(rng)
			encoded, err := marshalChoiceVariant(original, 3)
			if err != nil {
				t.Fatalf("iteration %d: marshal IdentificationBasedDetectionData failed: %v", i, err)
			}
			var decoded IdentificationBasedDetectionData
			err = unmarshalChoiceVariant(encoded, &decoded)
			if err != nil {
				t.Fatalf("iteration %d: unmarshal IdentificationBasedDetectionData failed: %v", i, err)
			}
			if !reflect.DeepEqual(original, decoded) {
				t.Fatalf("iteration %d: IdentificationBasedDetectionData round-trip mismatch\noriginal: %+v\ndecoded:  %+v", i, original, decoded)
			}

		case 4:
			original := genOccupancyBasedDetectionDataType2(rng)
			encoded, err := marshalChoiceVariant(original, 4)
			if err != nil {
				t.Fatalf("iteration %d: marshal OccupancyBasedDetectionDataType2 failed: %v", i, err)
			}
			var decoded OccupancyBasedDetectionDataType2
			err = unmarshalChoiceVariant(encoded, &decoded)
			if err != nil {
				t.Fatalf("iteration %d: unmarshal OccupancyBasedDetectionDataType2 failed: %v", i, err)
			}
			if !reflect.DeepEqual(original, decoded) {
				t.Fatalf("iteration %d: OccupancyBasedDetectionDataType2 round-trip mismatch\noriginal: %+v\ndecoded:  %+v", i, original, decoded)
			}

		case 5:
			original := genImageProcessingBasedDetectionDataType2(rng)
			encoded, err := marshalChoiceVariant(original, 5)
			if err != nil {
				t.Fatalf("iteration %d: marshal ImageProcessingBasedDetectionDataType2 failed: %v", i, err)
			}
			var decoded ImageProcessingBasedDetectionDataType2
			err = unmarshalChoiceVariant(encoded, &decoded)
			if err != nil {
				t.Fatalf("iteration %d: unmarshal ImageProcessingBasedDetectionDataType2 failed: %v", i, err)
			}
			if !reflect.DeepEqual(original, decoded) {
				t.Fatalf("iteration %d: ImageProcessingBasedDetectionDataType2 round-trip mismatch\noriginal: %+v\ndecoded:  %+v", i, original, decoded)
			}

		case 6:
			original := genIdentificationBasedDetectionDataType2(rng)
			encoded, err := marshalChoiceVariant(original, 6)
			if err != nil {
				t.Fatalf("iteration %d: marshal IdentificationBasedDetectionDataType2 failed: %v", i, err)
			}
			var decoded IdentificationBasedDetectionDataType2
			err = unmarshalChoiceVariant(encoded, &decoded)
			if err != nil {
				t.Fatalf("iteration %d: unmarshal IdentificationBasedDetectionDataType2 failed: %v", i, err)
			}
			if !reflect.DeepEqual(original, decoded) {
				t.Fatalf("iteration %d: IdentificationBasedDetectionDataType2 round-trip mismatch\noriginal: %+v\ndecoded:  %+v", i, original, decoded)
			}
		}
	}
}
