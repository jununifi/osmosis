// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: osmosis/downtimedetector/v1beta1/downtime_duration.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Downtime int32

const (
	Downtime_DURATION_30S  Downtime = 0
	Downtime_DURATION_1M   Downtime = 1
	Downtime_DURATION_2M   Downtime = 2
	Downtime_DURATION_3M   Downtime = 3
	Downtime_DURATION_4M   Downtime = 4
	Downtime_DURATION_5M   Downtime = 5
	Downtime_DURATION_10M  Downtime = 6
	Downtime_DURATION_20M  Downtime = 7
	Downtime_DURATION_30M  Downtime = 8
	Downtime_DURATION_40M  Downtime = 9
	Downtime_DURATION_50M  Downtime = 10
	Downtime_DURATION_1H   Downtime = 11
	Downtime_DURATION_1_5H Downtime = 12
	Downtime_DURATION_2H   Downtime = 13
	Downtime_DURATION_2_5H Downtime = 14
	Downtime_DURATION_3H   Downtime = 15
	Downtime_DURATION_4H   Downtime = 16
	Downtime_DURATION_5H   Downtime = 17
	Downtime_DURATION_6H   Downtime = 18
	Downtime_DURATION_9H   Downtime = 19
	Downtime_DURATION_12H  Downtime = 20
	Downtime_DURATION_18H  Downtime = 21
	Downtime_DURATION_24H  Downtime = 22
	Downtime_DURATION_36H  Downtime = 23
	Downtime_DURATION_48H  Downtime = 24
)

var Downtime_name = map[int32]string{
	0:  "DURATION_30S",
	1:  "DURATION_1M",
	2:  "DURATION_2M",
	3:  "DURATION_3M",
	4:  "DURATION_4M",
	5:  "DURATION_5M",
	6:  "DURATION_10M",
	7:  "DURATION_20M",
	8:  "DURATION_30M",
	9:  "DURATION_40M",
	10: "DURATION_50M",
	11: "DURATION_1H",
	12: "DURATION_1_5H",
	13: "DURATION_2H",
	14: "DURATION_2_5H",
	15: "DURATION_3H",
	16: "DURATION_4H",
	17: "DURATION_5H",
	18: "DURATION_6H",
	19: "DURATION_9H",
	20: "DURATION_12H",
	21: "DURATION_18H",
	22: "DURATION_24H",
	23: "DURATION_36H",
	24: "DURATION_48H",
}

var Downtime_value = map[string]int32{
	"DURATION_30S":  0,
	"DURATION_1M":   1,
	"DURATION_2M":   2,
	"DURATION_3M":   3,
	"DURATION_4M":   4,
	"DURATION_5M":   5,
	"DURATION_10M":  6,
	"DURATION_20M":  7,
	"DURATION_30M":  8,
	"DURATION_40M":  9,
	"DURATION_50M":  10,
	"DURATION_1H":   11,
	"DURATION_1_5H": 12,
	"DURATION_2H":   13,
	"DURATION_2_5H": 14,
	"DURATION_3H":   15,
	"DURATION_4H":   16,
	"DURATION_5H":   17,
	"DURATION_6H":   18,
	"DURATION_9H":   19,
	"DURATION_12H":  20,
	"DURATION_18H":  21,
	"DURATION_24H":  22,
	"DURATION_36H":  23,
	"DURATION_48H":  24,
}

func (x Downtime) String() string {
	return proto.EnumName(Downtime_name, int32(x))
}

func (Downtime) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0e46a3f710ad43a8, []int{0}
}

func init() {
	proto.RegisterEnum("osmosis.downtimedetector.v1beta1.Downtime", Downtime_name, Downtime_value)
}

func init() {
	proto.RegisterFile("osmosis/downtimedetector/v1beta1/downtime_duration.proto", fileDescriptor_0e46a3f710ad43a8)
}

var fileDescriptor_0e46a3f710ad43a8 = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x3d, 0x6f, 0xe2, 0x40,
	0x10, 0x86, 0xed, 0xe3, 0x8e, 0xe3, 0x0c, 0x1c, 0x83, 0x8f, 0xfb, 0xa2, 0xf0, 0x5d, 0x1d, 0x09,
	0xaf, 0x3f, 0x00, 0x91, 0x22, 0x45, 0x22, 0x8a, 0x4d, 0xb1, 0x89, 0x94, 0x0f, 0x45, 0x4a, 0x63,
	0xd9, 0xe0, 0x38, 0x96, 0x30, 0x8b, 0xf0, 0x42, 0xc2, 0xbf, 0xc8, 0x6f, 0x4a, 0x95, 0x92, 0x32,
	0x65, 0x04, 0x7f, 0x24, 0xc2, 0x36, 0x8e, 0xc6, 0x9d, 0xe7, 0x99, 0x79, 0xbd, 0xef, 0x3b, 0x1a,
	0x65, 0xc0, 0xe3, 0x88, 0xc7, 0x61, 0x4c, 0xc6, 0xfc, 0x61, 0x2a, 0xc2, 0xc8, 0x1f, 0xfb, 0xc2,
	0x1f, 0x09, 0x3e, 0x27, 0x4b, 0xd3, 0xf3, 0x85, 0x6b, 0xe6, 0x0d, 0x67, 0xbc, 0x98, 0xbb, 0x22,
	0xe4, 0x53, 0x7d, 0x36, 0xe7, 0x82, 0xab, 0xff, 0x33, 0xa5, 0x5e, 0x54, 0xea, 0x99, 0xb2, 0xdd,
	0x0a, 0x78, 0xc0, 0x93, 0x61, 0xb2, 0xfb, 0x4a, 0x75, 0xed, 0xbf, 0x01, 0xe7, 0xc1, 0xc4, 0x27,
	0x49, 0xe5, 0x2d, 0xee, 0x88, 0x3b, 0x5d, 0xed, 0x5b, 0xa3, 0xe4, 0x9f, 0x4e, 0xaa, 0x49, 0x8b,
	0xac, 0xa5, 0x15, 0x55, 0xd8, 0x4d, 0xfb, 0x5f, 0xb1, 0xbf, 0x73, 0x14, 0x0b, 0x37, 0x9a, 0xa5,
	0x03, 0x07, 0xcf, 0x25, 0xa5, 0x32, 0xcc, 0x9c, 0xaa, 0xa0, 0xd4, 0x86, 0xd7, 0x17, 0xc7, 0x57,
	0xa7, 0xe7, 0x67, 0x8e, 0x6d, 0x5c, 0x82, 0xa4, 0x36, 0x94, 0x6a, 0x4e, 0x4c, 0x06, 0x32, 0x02,
	0x16, 0x83, 0x4f, 0x08, 0xd8, 0x0c, 0x4a, 0x08, 0x74, 0x19, 0x7c, 0x46, 0xa0, 0xc7, 0xe0, 0x0b,
	0x7a, 0xc6, 0x34, 0x18, 0x94, 0x11, 0xb1, 0x0c, 0x06, 0x5f, 0x0b, 0x56, 0x18, 0x54, 0x10, 0xe9,
	0x1a, 0x0c, 0xbe, 0x21, 0xd2, 0x33, 0x18, 0x28, 0xd8, 0x2e, 0x85, 0xaa, 0xda, 0x54, 0xea, 0x1f,
	0xc0, 0xe9, 0x51, 0xa8, 0xe1, 0x04, 0x14, 0xea, 0x68, 0xc6, 0xda, 0xcd, 0x7c, 0xc7, 0xa1, 0x28,
	0x34, 0x70, 0x28, 0x0a, 0x80, 0x43, 0x51, 0x68, 0x22, 0xd0, 0xa7, 0xa0, 0x22, 0x70, 0x48, 0xe1,
	0x07, 0x8e, 0x6d, 0x51, 0x68, 0x61, 0x32, 0xa0, 0xf0, 0x13, 0x2f, 0xa2, 0x4b, 0xe1, 0x17, 0x5e,
	0x44, 0x9f, 0xc2, 0x6f, 0xbc, 0x88, 0x01, 0x85, 0x3f, 0x27, 0x37, 0x2f, 0x1b, 0x4d, 0x5e, 0x6f,
	0x34, 0xf9, 0x6d, 0xa3, 0xc9, 0x4f, 0x5b, 0x4d, 0x5a, 0x6f, 0x35, 0xe9, 0x75, 0xab, 0x49, 0xb7,
	0x47, 0x41, 0x28, 0xee, 0x17, 0x9e, 0x3e, 0xe2, 0x11, 0xc9, 0x0e, 0xb3, 0x33, 0x71, 0xbd, 0x78,
	0x5f, 0x90, 0xa5, 0x65, 0x93, 0xc7, 0xfc, 0x98, 0x3b, 0xf9, 0x99, 0x8b, 0xd5, 0xcc, 0x8f, 0xbd,
	0x72, 0x72, 0x24, 0xf6, 0x7b, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x2e, 0xa0, 0x2e, 0x0f, 0x03,
	0x00, 0x00,
}
