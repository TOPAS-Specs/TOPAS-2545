"""
Property-based test for Python BER encoding round-trip.

Feature: asn1-spec-update, Property 2: Python BER encoding round-trip

For any valid IpmstscdData structure decoded from JER, encoding to BER
and decoding back should produce a structurally equal value.

Validates: Requirements 7.1
"""

import os
import sys

import asn1tools
from hypothesis import given, settings, HealthCheck
from hypothesis import strategies as st

# Compile the ASN.1 schema for JER and BER codecs
ASN_FILE = os.path.join(os.path.dirname(__file__), '..', 'Annex_A_10711_v3.0.asn')
jer_codec = asn1tools.compile_files(ASN_FILE, 'jer')
ber_codec = asn1tools.compile_files(ASN_FILE, 'ber')

TYPE_NAME = 'IpmstscdData'


# --- Hypothesis strategies for generating valid IpmstscdData structures ---

@st.composite
def time_strategy(draw):
    """Generate a valid Time structure."""
    time_val = {}
    if draw(st.booleans()):
        time_val['year'] = draw(st.integers(min_value=-32768, max_value=32767))
    if draw(st.booleans()):
        time_val['month'] = draw(st.integers(min_value=1, max_value=12))
    if draw(st.booleans()):
        time_val['day'] = draw(st.integers(min_value=1, max_value=31))
    time_val['hour'] = draw(st.integers(min_value=0, max_value=23))
    time_val['minute'] = draw(st.integers(min_value=0, max_value=59))
    time_val['second'] = draw(st.integers(min_value=0, max_value=60))
    frac_choice = draw(st.sampled_from(['deciSeconds', 'centiSeconds', 'milliseconds']))
    if frac_choice == 'deciSeconds':
        time_val['secondFractions'] = ('deciSeconds', draw(st.integers(min_value=0, max_value=9)))
    elif frac_choice == 'centiSeconds':
        time_val['secondFractions'] = ('centiSeconds', draw(st.integers(min_value=0, max_value=99)))
    else:
        time_val['secondFractions'] = ('milliseconds', draw(st.integers(min_value=0, max_value=999)))
    if draw(st.booleans()):
        time_val['timeZone'] = {
            'hour': draw(st.integers(min_value=-13, max_value=13)),
            'minute': draw(st.integers(min_value=0, max_value=59)),
        }
    return time_val


@st.composite
def general_time_location_core_strategy(draw):
    """Generate a valid GeneralTimeLocationCore structure."""
    val = {
        'currentTime': draw(time_strategy()),
    }
    if draw(st.booleans()):
        val['locationLongitude'] = draw(st.integers(min_value=-180000000, max_value=180000000))
    if draw(st.booleans()):
        val['locationLatitude'] = draw(st.integers(min_value=-90000000, max_value=90000000))
    if draw(st.booleans()):
        val['locationElevation'] = draw(st.integers(min_value=-8192, max_value=57344))
    return val


@st.composite
def occupancy_non_occupancy_history_strategy(draw):
    """Generate a valid OccupancyNonOccupancyHistory entry."""
    return {
        'occupancyTimes': draw(st.integers(min_value=0, max_value=65535)),
        'nonOccupancyTimes': draw(st.integers(min_value=0, max_value=65535)),
    }


@st.composite
def occupancy_based_detection_data_strategy(draw):
    """Generate a valid OccupancyBasedDetectionData structure."""
    val = {}
    if draw(st.booleans()):
        val['occupancyDataDuration'] = draw(st.integers(min_value=0, max_value=65535))
    val['occupancyState'] = draw(st.booleans())
    val['occupancyStateDuration'] = draw(st.integers(min_value=0, max_value=65535))
    val['occupancyPreviousStateDuration'] = draw(st.integers(min_value=0, max_value=65535))
    val['occupancyRate'] = draw(st.floats(min_value=-1000.0, max_value=1000.0, allow_nan=False, allow_infinity=False))
    if draw(st.booleans()):
        val['speed'] = draw(st.floats(min_value=0.0, max_value=300.0, allow_nan=False, allow_infinity=False))
    val['volume'] = draw(st.integers(min_value=0, max_value=65535))
    if draw(st.booleans()):
        val['occupancyNonOccupancyHistory'] = draw(
            st.lists(occupancy_non_occupancy_history_strategy(), min_size=1, max_size=3)
        )
    if draw(st.booleans()):
        val['errorState'] = draw(st.sampled_from([
            'openCircuit', 'shortCircuit', 'occupancyError',
            'nonOccupancyError', 'volumeError', 'parameterInvalid',
            'managementNeeded'
        ]))
    if draw(st.booleans()):
        val['userData'] = draw(st.binary(min_size=1, max_size=8))
    if draw(st.booleans()):
        val['targetType'] = draw(st.integers(min_value=1, max_value=255))
    if draw(st.booleans()):
        val['directionDiscrimination'] = draw(st.booleans())
    return val


@st.composite
def image_processing_based_detection_data_strategy(draw):
    """Generate a valid ImageProcessingBasedDetectionData structure."""
    val = {}
    if draw(st.booleans()):
        val['dataDuration'] = draw(st.integers(min_value=0, max_value=65535))
    val['volumeQuantity'] = draw(st.integers(min_value=0, max_value=65535))
    if draw(st.booleans()):
        val['queueLength'] = draw(st.integers(min_value=0, max_value=65535))
    if draw(st.booleans()):
        val['occupancyRate'] = draw(st.floats(min_value=0.0, max_value=100.0, allow_nan=False, allow_infinity=False))
    if draw(st.booleans()):
        val['speedQuantity'] = draw(st.floats(min_value=0.0, max_value=300.0, allow_nan=False, allow_infinity=False))
    if draw(st.booleans()):
        val['occupancyNonOccupancyHistory'] = draw(occupancy_non_occupancy_history_strategy())
    if draw(st.booleans()):
        val['errorState'] = draw(st.sampled_from([
            'deviceFail', 'unstableUtility', 'connectionFail',
            'imageProcessingFail', 'parameterInvalid', 'volumeError',
            'managementNeeded'
        ]))
    if draw(st.booleans()):
        val['userData'] = draw(st.binary(min_size=1, max_size=8))
    return val


@st.composite
def identification_based_detection_data_strategy(draw):
    """Generate a valid IdentificationBasedDetectionData structure."""
    val = {
        'sequenceNumber': draw(st.integers(min_value=0, max_value=255)),
    }
    if draw(st.booleans()):
        val['deviceType'] = draw(st.sampled_from([
            'infrared', 'radioFrequency', 'vehicleDetectionSystems',
            'magnetics', 'barCodeScanner', 'tagScanner', 'other'
        ]))
    val['vehicleId'] = draw(st.binary(min_size=1, max_size=16))
    if draw(st.booleans()):
        val['vehicleType'] = draw(st.integers(min_value=0, max_value=255))
    if draw(st.booleans()):
        val['vehicleUse'] = draw(st.integers(min_value=0, max_value=255))
    if draw(st.booleans()):
        val['detectionLane'] = draw(st.integers(min_value=1, max_value=8))
    if draw(st.booleans()):
        val['detectionLaneMedian'] = draw(st.integers(min_value=1, max_value=8))
    if draw(st.booleans()):
        val['detectionSpeed'] = draw(st.floats(min_value=0.0, max_value=300.0, allow_nan=False, allow_infinity=False))
    if draw(st.booleans()):
        val['occupancy'] = draw(st.integers(min_value=0, max_value=65535))
    if draw(st.booleans()):
        val['errorState'] = draw(st.sampled_from([
            'rSEFail', 'rSEConnectionFail', 'wirelessFail',
            'unstableUtility', 'managementNeeded'
        ]))
    if draw(st.booleans()):
        val['tagInfo'] = draw(st.binary(min_size=1, max_size=16))
    if draw(st.booleans()):
        val['userData'] = draw(st.binary(min_size=1, max_size=8))
    return val


@st.composite
def accumulative_detection_seq_strategy(draw):
    """Generate a valid AccumulativeDetectionSeq entry."""
    val = {
        'detectorNumber': draw(st.integers(min_value=1, max_value=48)),
    }
    if draw(st.booleans()):
        val['detectionStatus'] = draw(st.sampled_from(['normal', 'failure', 'dataInvalid']))
    val['density'] = draw(st.integers(min_value=0, max_value=65535))
    val['occupancy'] = draw(st.integers(min_value=0, max_value=65535))
    val['detectorPulseError'] = draw(st.integers(min_value=0, max_value=65535))
    return val


@st.composite
def passing_vehicle_time_series_seq_strategy(draw):
    """Generate a valid PassingVehicleTimeSeriesSeq entry."""
    val = {
        'detectorNumber': draw(st.integers(min_value=1, max_value=48)),
    }
    if draw(st.booleans()):
        val['detectionStatus'] = draw(st.sampled_from(['normal', 'failure', 'dataInvalid']))
    val['vehiclePassageHistory'] = draw(st.binary(min_size=8, max_size=8))
    return val


@st.composite
def vehicle_speed_seq_strategy(draw):
    """Generate a valid VehicleSpeedSeq entry."""
    return {
        'detectorNumber': draw(st.integers(min_value=1, max_value=48)),
        'vehicleType': draw(st.sampled_from([
            'fourBus', 'fourLargeSizeTruck', 'fourSmallSizeTruck',
            'fourElse', 'twoLargeSizeVehicle', 'twoElse'
        ])),
        'velocity': draw(st.integers(min_value=0, max_value=127)),
    }


@st.composite
def vehicle_identification_seq_strategy(draw):
    """Generate a valid VehicleIdentificationSeq entry."""
    val = {
        'vehicleID': draw(st.binary(min_size=1, max_size=16)),
    }
    if draw(st.booleans()):
        val['data'] = draw(st.binary(min_size=1, max_size=16))
    return val


@st.composite
def detector_status_strategy(draw):
    """Generate a valid DetectorStatus BIT STRING (1..8 bits).

    asn1tools represents BIT STRING as tuple(bytes, number_of_bits).
    DetectorStatus has named bits: processingFailure(0), operatingFailure(1).
    SIZE constraint is 1..8 bits.
    """
    # Generate a byte with up to 2 meaningful bits (positions 0 and 1)
    bit0 = draw(st.booleans())  # processingFailure
    bit1 = draw(st.booleans())  # operatingFailure
    byte_val = (0x80 if bit0 else 0) | (0x40 if bit1 else 0)
    return (bytes([byte_val]), 2)


@st.composite
def directional_traffic_volume_seq_strategy(draw):
    """Generate a valid DirectionalTrafficVolumeSeq entry."""
    val = {
        'directionNumber': draw(st.integers(min_value=1, max_value=32)),
    }
    if draw(st.booleans()):
        val['detectionStatus'] = draw(st.sampled_from(['normal', 'failure', 'dataInvalid']))
    val['directionDensity'] = draw(st.integers(min_value=0, max_value=65535))
    return val


@st.composite
def occupancy_based_detection_data_type2_strategy(draw):
    """Generate a valid OccupancyBasedDetectionDataType2 structure."""
    val = {
        'accumulativeDetection': draw(
            st.lists(accumulative_detection_seq_strategy(), min_size=1, max_size=3)
        ),
        'passingVehicleTimeSeries': draw(
            st.lists(passing_vehicle_time_series_seq_strategy(), min_size=1, max_size=3)
        ),
        'vehicleSpeed': draw(
            st.lists(vehicle_speed_seq_strategy(), min_size=0, max_size=3)
        ),
        'specificVehicleDetection': draw(st.binary(min_size=6, max_size=6)),
    }
    if draw(st.booleans()):
        val['detectorStatus'] = draw(detector_status_strategy())
    return val


@st.composite
def image_processing_based_detection_data_type2_strategy(draw):
    """Generate a valid ImageProcessingBasedDetectionDataType2 structure."""
    val = {
        'accumulativeDetection': draw(
            st.lists(accumulative_detection_seq_strategy(), min_size=1, max_size=3)
        ),
        'passingVehicleTimeSeries': draw(
            st.lists(passing_vehicle_time_series_seq_strategy(), min_size=1, max_size=3)
        ),
        'vehicleSpeed': draw(
            st.lists(vehicle_speed_seq_strategy(), min_size=0, max_size=3)
        ),
        'specificVehicleDetection': draw(st.binary(min_size=6, max_size=6)),
        'congestionInfo': {
            'congestionLength1': draw(st.integers(min_value=0, max_value=150)),
            'vehicleStartPosition1': draw(st.integers(min_value=0, max_value=150)),
            'congestionLength2': draw(st.integers(min_value=0, max_value=150)),
            'vehicleStartPosition2': draw(st.integers(min_value=0, max_value=150)),
        },
        'directionalTrafficVolume': draw(
            st.lists(directional_traffic_volume_seq_strategy(), min_size=1, max_size=3)
        ),
    }
    if draw(st.booleans()):
        val['detectorStatus'] = draw(detector_status_strategy())
    return val


@st.composite
def identification_based_detection_data_type2_strategy(draw):
    """Generate a valid IdentificationBasedDetectionDataType2 structure."""
    val = {
        'accumulativeDetection': draw(
            st.lists(accumulative_detection_seq_strategy(), min_size=1, max_size=3)
        ),
        'passingVehicleTimeSeries': draw(
            st.lists(passing_vehicle_time_series_seq_strategy(), min_size=1, max_size=3)
        ),
        'vehicleSpeed': draw(
            st.lists(vehicle_speed_seq_strategy(), min_size=0, max_size=3)
        ),
        'specificVehicleDetection': draw(st.binary(min_size=6, max_size=6)),
        'vehicleIdentification': draw(
            st.lists(vehicle_identification_seq_strategy(), min_size=1, max_size=3)
        ),
    }
    if draw(st.booleans()):
        val['detectorStatus'] = draw(detector_status_strategy())
    return val


@st.composite
def detection_info_strategy(draw):
    """Generate a valid detectionInfo CHOICE value (one of 6 variants)."""
    variant = draw(st.sampled_from([
        'occupancyBasedDetectionData',
        'imageProcessingBasedDetectionData',
        'identificationBasedDetectionData',
        'occupancyTypeDetectionInfoType2',
        'imageProcessingBasedDetectionDataType2',
        'identificationBasedDetectionDataType2',
    ]))
    if variant == 'occupancyBasedDetectionData':
        return (variant, draw(occupancy_based_detection_data_strategy()))
    elif variant == 'imageProcessingBasedDetectionData':
        return (variant, draw(image_processing_based_detection_data_strategy()))
    elif variant == 'identificationBasedDetectionData':
        return (variant, draw(identification_based_detection_data_strategy()))
    elif variant == 'occupancyTypeDetectionInfoType2':
        return (variant, draw(occupancy_based_detection_data_type2_strategy()))
    elif variant == 'imageProcessingBasedDetectionDataType2':
        return (variant, draw(image_processing_based_detection_data_type2_strategy()))
    else:
        return (variant, draw(identification_based_detection_data_type2_strategy()))


# Map CHOICE variant names to their corresponding infoType enum values
CHOICE_TO_INFO_TYPE = {
    'occupancyBasedDetectionData': 'occupancyTypeDetector',
    'imageProcessingBasedDetectionData': 'imageTypeDetector',
    'identificationBasedDetectionData': 'idBaseTypeDetector',
    'occupancyTypeDetectionInfoType2': 'occupancyTypeDetector',
    'imageProcessingBasedDetectionDataType2': 'imageTypeDetector',
    'identificationBasedDetectionDataType2': 'idBaseTypeDetector',
}


@st.composite
def detector_info_seq_strategy(draw):
    """Generate a valid DetectorInfoSeq with any detection data variant."""
    detection_info = draw(detection_info_strategy())
    choice_name = detection_info[0]
    val = {
        'physicalDetectorIndex': draw(st.integers(min_value=0, max_value=255)),
        'infoType': CHOICE_TO_INFO_TYPE[choice_name],
        'detectionInfo': detection_info,
    }
    if draw(st.booleans()):
        val['detectorTimeLocation'] = draw(general_time_location_core_strategy())
    return val


@st.composite
def detector_controller_info_strategy(draw):
    """Generate a valid DetectorControllerInfo structure."""
    val = {
        'detectorControllerIndex': draw(st.integers(min_value=0, max_value=255)),
    }
    if draw(st.booleans()):
        val['detectorControllerTimeLocation'] = draw(general_time_location_core_strategy())
    return val


@st.composite
def ipmstscd_data_strategy(draw):
    """Generate a valid IpmstscdData structure with random detection data variants."""
    val = {}
    if draw(st.booleans()):
        val['detectorControllerInfo'] = draw(detector_controller_info_strategy())
    val['detectorInfo'] = draw(
        st.lists(detector_info_seq_strategy(), min_size=1, max_size=4)
    )
    return val


# --- Property test ---

@settings(max_examples=100, suppress_health_check=[HealthCheck.too_slow])
@given(data=ipmstscd_data_strategy())
def test_ber_roundtrip(data):
    """
    **Validates: Requirements 7.1**

    Property 2: Python BER encoding round-trip

    For any valid IpmstscdData structure, encoding to BER and decoding back
    should produce a structurally equal value.
    """
    # Encode to BER
    ber_bytes = ber_codec.encode(TYPE_NAME, data)

    # Decode from BER
    decoded = ber_codec.decode(TYPE_NAME, ber_bytes)

    # Assert structural equality
    assert decoded == data, (
        f"BER round-trip failed.\n"
        f"Original: {data}\n"
        f"Decoded:  {decoded}"
    )


if __name__ == '__main__':
    test_ber_roundtrip()
    print("All property tests passed!")
