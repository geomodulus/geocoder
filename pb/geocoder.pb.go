// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: geocoder.proto

package pb

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

type LngLat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lng float64 `protobuf:"fixed64,1,opt,name=lng,proto3" json:"lng,omitempty"`
	Lat float64 `protobuf:"fixed64,2,opt,name=lat,proto3" json:"lat,omitempty"`
}

func (x *LngLat) Reset() {
	*x = LngLat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_geocoder_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LngLat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LngLat) ProtoMessage() {}

func (x *LngLat) ProtoReflect() protoreflect.Message {
	mi := &file_geocoder_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LngLat.ProtoReflect.Descriptor instead.
func (*LngLat) Descriptor() ([]byte, []int) {
	return file_geocoder_proto_rawDescGZIP(), []int{0}
}

func (x *LngLat) GetLng() float64 {
	if x != nil {
		return x.Lng
	}
	return 0
}

func (x *LngLat) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number string `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Street string `protobuf:"bytes,2,opt,name=street,proto3" json:"street,omitempty"`
	// Cross street, when present, indicates this is an intersection.
	CrossStreet string  `protobuf:"bytes,3,opt,name=cross_street,json=crossStreet,proto3" json:"cross_street,omitempty"`
	Location    *LngLat `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	Desc        string  `protobuf:"bytes,5,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_geocoder_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_geocoder_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_geocoder_proto_rawDescGZIP(), []int{1}
}

func (x *Location) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *Location) GetStreet() string {
	if x != nil {
		return x.Street
	}
	return ""
}

func (x *Location) GetCrossStreet() string {
	if x != nil {
		return x.CrossStreet
	}
	return ""
}

func (x *Location) GetLocation() *LngLat {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *Location) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

var File_geocoder_proto protoreflect.FileDescriptor

var file_geocoder_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x67, 0x65, 0x6f, 0x63, 0x6f, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2c, 0x0a, 0x06, 0x4c, 0x6e, 0x67, 0x4c, 0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6e,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03,
	0x6c, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x61, 0x74, 0x22, 0x96,
	0x01, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x63,
	0x72, 0x6f, 0x73, 0x73, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x12, 0x23,
	0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x07, 0x2e, 0x4c, 0x6e, 0x67, 0x4c, 0x61, 0x74, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_geocoder_proto_rawDescOnce sync.Once
	file_geocoder_proto_rawDescData = file_geocoder_proto_rawDesc
)

func file_geocoder_proto_rawDescGZIP() []byte {
	file_geocoder_proto_rawDescOnce.Do(func() {
		file_geocoder_proto_rawDescData = protoimpl.X.CompressGZIP(file_geocoder_proto_rawDescData)
	})
	return file_geocoder_proto_rawDescData
}

var file_geocoder_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_geocoder_proto_goTypes = []interface{}{
	(*LngLat)(nil),   // 0: LngLat
	(*Location)(nil), // 1: Location
}
var file_geocoder_proto_depIdxs = []int32{
	0, // 0: Location.location:type_name -> LngLat
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_geocoder_proto_init() }
func file_geocoder_proto_init() {
	if File_geocoder_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_geocoder_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LngLat); i {
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
		file_geocoder_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_geocoder_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_geocoder_proto_goTypes,
		DependencyIndexes: file_geocoder_proto_depIdxs,
		MessageInfos:      file_geocoder_proto_msgTypes,
	}.Build()
	File_geocoder_proto = out.File
	file_geocoder_proto_rawDesc = nil
	file_geocoder_proto_goTypes = nil
	file_geocoder_proto_depIdxs = nil
}
