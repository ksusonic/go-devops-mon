// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: proto/metrics/metric.proto

package metrics

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetricType int32

const (
	MetricType_unknown MetricType = 0
	MetricType_gauge   MetricType = 1
	MetricType_counter MetricType = 2
)

// Enum value maps for MetricType.
var (
	MetricType_name = map[int32]string{
		0: "unknown",
		1: "gauge",
		2: "counter",
	}
	MetricType_value = map[string]int32{
		"unknown": 0,
		"gauge":   1,
		"counter": 2,
	}
)

func (x MetricType) Enum() *MetricType {
	p := new(MetricType)
	*p = x
	return p
}

func (x MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_metrics_metric_proto_enumTypes[0].Descriptor()
}

func (MetricType) Type() protoreflect.EnumType {
	return &file_proto_metrics_metric_proto_enumTypes[0]
}

func (x MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetricType.Descriptor instead.
func (MetricType) EnumDescriptor() ([]byte, []int) {
	return file_proto_metrics_metric_proto_rawDescGZIP(), []int{0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID   string     `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Type MetricType `protobuf:"varint,2,opt,name=Type,proto3,enum=metrics.MetricType" json:"Type,omitempty"`
	// Types that are assignable to Payload:
	//
	//	*Metric_Delta
	//	*Metric_Value
	Payload isMetric_Payload `protobuf_oneof:"Payload"`
	Hash    string           `protobuf:"bytes,5,opt,name=Hash,proto3" json:"Hash,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrics_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metrics_metric_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Metric) GetType() MetricType {
	if x != nil {
		return x.Type
	}
	return MetricType_unknown
}

func (m *Metric) GetPayload() isMetric_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *Metric) GetDelta() int64 {
	if x, ok := x.GetPayload().(*Metric_Delta); ok {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x, ok := x.GetPayload().(*Metric_Value); ok {
		return x.Value
	}
	return 0
}

func (x *Metric) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type isMetric_Payload interface {
	isMetric_Payload()
}

type Metric_Delta struct {
	Delta int64 `protobuf:"varint,3,opt,name=Delta,proto3,oneof"`
}

type Metric_Value struct {
	Value float64 `protobuf:"fixed64,4,opt,name=Value,proto3,oneof"`
}

func (*Metric_Delta) isMetric_Payload() {}

func (*Metric_Value) isMetric_Payload() {}

var File_proto_metrics_metric_proto protoreflect.FileDescriptor

var file_proto_metrics_metric_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2f,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x90, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x27, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13,
	0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x05, 0x44, 0x65, 0x6c,
	0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x74,
	0x61, 0x12, 0x16, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01,
	0x48, 0x00, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x42, 0x09, 0x0a,
	0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2a, 0x31, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77,
	0x6e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x67, 0x61, 0x75, 0x67, 0x65, 0x10, 0x01, 0x12, 0x0b,
	0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x10, 0x02, 0x42, 0x31, 0x5a, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x73, 0x75, 0x73, 0x6f, 0x6e,
	0x69, 0x63, 0x2f, 0x67, 0x6f, 0x2d, 0x64, 0x65, 0x76, 0x6f, 0x70, 0x73, 0x2d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metrics_metric_proto_rawDescOnce sync.Once
	file_proto_metrics_metric_proto_rawDescData = file_proto_metrics_metric_proto_rawDesc
)

func file_proto_metrics_metric_proto_rawDescGZIP() []byte {
	file_proto_metrics_metric_proto_rawDescOnce.Do(func() {
		file_proto_metrics_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metrics_metric_proto_rawDescData)
	})
	return file_proto_metrics_metric_proto_rawDescData
}

var file_proto_metrics_metric_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_metrics_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_metrics_metric_proto_goTypes = []interface{}{
	(MetricType)(0), // 0: metrics.MetricType
	(*Metric)(nil),  // 1: metrics.Metric
}
var file_proto_metrics_metric_proto_depIdxs = []int32{
	0, // 0: metrics.Metric.Type:type_name -> metrics.MetricType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_metrics_metric_proto_init() }
func file_proto_metrics_metric_proto_init() {
	if File_proto_metrics_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metrics_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_metrics_metric_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Metric_Delta)(nil),
		(*Metric_Value)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metrics_metric_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_metrics_metric_proto_goTypes,
		DependencyIndexes: file_proto_metrics_metric_proto_depIdxs,
		EnumInfos:         file_proto_metrics_metric_proto_enumTypes,
		MessageInfos:      file_proto_metrics_metric_proto_msgTypes,
	}.Build()
	File_proto_metrics_metric_proto = out.File
	file_proto_metrics_metric_proto_rawDesc = nil
	file_proto_metrics_metric_proto_goTypes = nil
	file_proto_metrics_metric_proto_depIdxs = nil
}