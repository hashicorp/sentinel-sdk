// Code generated by protoc-gen-go. DO NOT EDIT.
// source: import.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Type is an enum representing the type of the value. This isn't the
// full set of Sentinel types since some types cannot be sent via
// Protobufs such as rules or functions.
type Value_Type int32

const (
	Value_INVALID   Value_Type = 0
	Value_UNDEFINED Value_Type = 1
	Value_NULL      Value_Type = 2
	Value_BOOL      Value_Type = 3
	Value_INT       Value_Type = 4
	Value_FLOAT     Value_Type = 5
	Value_STRING    Value_Type = 6
	Value_LIST      Value_Type = 7
	Value_MAP       Value_Type = 8
)

var Value_Type_name = map[int32]string{
	0: "INVALID",
	1: "UNDEFINED",
	2: "NULL",
	3: "BOOL",
	4: "INT",
	5: "FLOAT",
	6: "STRING",
	7: "LIST",
	8: "MAP",
}
var Value_Type_value = map[string]int32{
	"INVALID":   0,
	"UNDEFINED": 1,
	"NULL":      2,
	"BOOL":      3,
	"INT":       4,
	"FLOAT":     5,
	"STRING":    6,
	"LIST":      7,
	"MAP":       8,
}

func (x Value_Type) String() string {
	return proto.EnumName(Value_Type_name, int32(x))
}
func (Value_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{4, 0}
}

// Empty is just an empty message.
type EmptyResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyResp) Reset()         { *m = EmptyResp{} }
func (m *EmptyResp) String() string { return proto.CompactTextString(m) }
func (*EmptyResp) ProtoMessage()    {}
func (*EmptyResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{0}
}
func (m *EmptyResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyResp.Unmarshal(m, b)
}
func (m *EmptyResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyResp.Marshal(b, m, deterministic)
}
func (dst *EmptyResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyResp.Merge(dst, src)
}
func (m *EmptyResp) XXX_Size() int {
	return xxx_messageInfo_EmptyResp.Size(m)
}
func (m *EmptyResp) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyResp.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyResp proto.InternalMessageInfo

// Configure are the structures for Import.Configure
type Configure struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Configure) Reset()         { *m = Configure{} }
func (m *Configure) String() string { return proto.CompactTextString(m) }
func (*Configure) ProtoMessage()    {}
func (*Configure) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{1}
}
func (m *Configure) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Configure.Unmarshal(m, b)
}
func (m *Configure) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Configure.Marshal(b, m, deterministic)
}
func (dst *Configure) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Configure.Merge(dst, src)
}
func (m *Configure) XXX_Size() int {
	return xxx_messageInfo_Configure.Size(m)
}
func (m *Configure) XXX_DiscardUnknown() {
	xxx_messageInfo_Configure.DiscardUnknown(m)
}

var xxx_messageInfo_Configure proto.InternalMessageInfo

type Configure_Request struct {
	Config               *Value   `protobuf:"bytes,3,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Configure_Request) Reset()         { *m = Configure_Request{} }
func (m *Configure_Request) String() string { return proto.CompactTextString(m) }
func (*Configure_Request) ProtoMessage()    {}
func (*Configure_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{1, 0}
}
func (m *Configure_Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Configure_Request.Unmarshal(m, b)
}
func (m *Configure_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Configure_Request.Marshal(b, m, deterministic)
}
func (dst *Configure_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Configure_Request.Merge(dst, src)
}
func (m *Configure_Request) XXX_Size() int {
	return xxx_messageInfo_Configure_Request.Size(m)
}
func (m *Configure_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Configure_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Configure_Request proto.InternalMessageInfo

func (m *Configure_Request) GetConfig() *Value {
	if m != nil {
		return m.Config
	}
	return nil
}

type Configure_Response struct {
	InstanceId           uint64   `protobuf:"varint,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Configure_Response) Reset()         { *m = Configure_Response{} }
func (m *Configure_Response) String() string { return proto.CompactTextString(m) }
func (*Configure_Response) ProtoMessage()    {}
func (*Configure_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{1, 1}
}
func (m *Configure_Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Configure_Response.Unmarshal(m, b)
}
func (m *Configure_Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Configure_Response.Marshal(b, m, deterministic)
}
func (dst *Configure_Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Configure_Response.Merge(dst, src)
}
func (m *Configure_Response) XXX_Size() int {
	return xxx_messageInfo_Configure_Response.Size(m)
}
func (m *Configure_Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Configure_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Configure_Response proto.InternalMessageInfo

func (m *Configure_Response) GetInstanceId() uint64 {
	if m != nil {
		return m.InstanceId
	}
	return 0
}

// Get are the structures for an Import.Get.
type Get struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Get) Reset()         { *m = Get{} }
func (m *Get) String() string { return proto.CompactTextString(m) }
func (*Get) ProtoMessage()    {}
func (*Get) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{2}
}
func (m *Get) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Get.Unmarshal(m, b)
}
func (m *Get) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Get.Marshal(b, m, deterministic)
}
func (dst *Get) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Get.Merge(dst, src)
}
func (m *Get) XXX_Size() int {
	return xxx_messageInfo_Get.Size(m)
}
func (m *Get) XXX_DiscardUnknown() {
	xxx_messageInfo_Get.DiscardUnknown(m)
}

var xxx_messageInfo_Get proto.InternalMessageInfo

// Request is a single request for a Get.
type Get_Request struct {
	InstanceId           uint64   `protobuf:"varint,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	ExecId               uint64   `protobuf:"varint,2,opt,name=exec_id,json=execId,proto3" json:"exec_id,omitempty"`
	ExecDeadline         uint64   `protobuf:"varint,3,opt,name=exec_deadline,json=execDeadline,proto3" json:"exec_deadline,omitempty"`
	Keys                 []string `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
	KeyId                uint64   `protobuf:"varint,5,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	Call                 bool     `protobuf:"varint,6,opt,name=call,proto3" json:"call,omitempty"`
	Args                 []*Value `protobuf:"bytes,7,rep,name=args,proto3" json:"args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Get_Request) Reset()         { *m = Get_Request{} }
func (m *Get_Request) String() string { return proto.CompactTextString(m) }
func (*Get_Request) ProtoMessage()    {}
func (*Get_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{2, 0}
}
func (m *Get_Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Get_Request.Unmarshal(m, b)
}
func (m *Get_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Get_Request.Marshal(b, m, deterministic)
}
func (dst *Get_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Get_Request.Merge(dst, src)
}
func (m *Get_Request) XXX_Size() int {
	return xxx_messageInfo_Get_Request.Size(m)
}
func (m *Get_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Get_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Get_Request proto.InternalMessageInfo

func (m *Get_Request) GetInstanceId() uint64 {
	if m != nil {
		return m.InstanceId
	}
	return 0
}

func (m *Get_Request) GetExecId() uint64 {
	if m != nil {
		return m.ExecId
	}
	return 0
}

func (m *Get_Request) GetExecDeadline() uint64 {
	if m != nil {
		return m.ExecDeadline
	}
	return 0
}

func (m *Get_Request) GetKeys() []string {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *Get_Request) GetKeyId() uint64 {
	if m != nil {
		return m.KeyId
	}
	return 0
}

func (m *Get_Request) GetCall() bool {
	if m != nil {
		return m.Call
	}
	return false
}

func (m *Get_Request) GetArgs() []*Value {
	if m != nil {
		return m.Args
	}
	return nil
}

// Response is a single response for a Get.
type Get_Response struct {
	InstanceId           uint64   `protobuf:"varint,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	KeyId                uint64   `protobuf:"varint,2,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	Keys                 []string `protobuf:"bytes,3,rep,name=keys,proto3" json:"keys,omitempty"`
	Value                *Value   `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Get_Response) Reset()         { *m = Get_Response{} }
func (m *Get_Response) String() string { return proto.CompactTextString(m) }
func (*Get_Response) ProtoMessage()    {}
func (*Get_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{2, 1}
}
func (m *Get_Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Get_Response.Unmarshal(m, b)
}
func (m *Get_Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Get_Response.Marshal(b, m, deterministic)
}
func (dst *Get_Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Get_Response.Merge(dst, src)
}
func (m *Get_Response) XXX_Size() int {
	return xxx_messageInfo_Get_Response.Size(m)
}
func (m *Get_Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Get_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Get_Response proto.InternalMessageInfo

func (m *Get_Response) GetInstanceId() uint64 {
	if m != nil {
		return m.InstanceId
	}
	return 0
}

func (m *Get_Response) GetKeyId() uint64 {
	if m != nil {
		return m.KeyId
	}
	return 0
}

func (m *Get_Response) GetKeys() []string {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *Get_Response) GetValue() *Value {
	if m != nil {
		return m.Value
	}
	return nil
}

// MultiRequest allows multiple requests in a single Get.
type Get_MultiRequest struct {
	Requests             []*Get_Request `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Get_MultiRequest) Reset()         { *m = Get_MultiRequest{} }
func (m *Get_MultiRequest) String() string { return proto.CompactTextString(m) }
func (*Get_MultiRequest) ProtoMessage()    {}
func (*Get_MultiRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{2, 2}
}
func (m *Get_MultiRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Get_MultiRequest.Unmarshal(m, b)
}
func (m *Get_MultiRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Get_MultiRequest.Marshal(b, m, deterministic)
}
func (dst *Get_MultiRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Get_MultiRequest.Merge(dst, src)
}
func (m *Get_MultiRequest) XXX_Size() int {
	return xxx_messageInfo_Get_MultiRequest.Size(m)
}
func (m *Get_MultiRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Get_MultiRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Get_MultiRequest proto.InternalMessageInfo

func (m *Get_MultiRequest) GetRequests() []*Get_Request {
	if m != nil {
		return m.Requests
	}
	return nil
}

// MultiResponse allows multiple responses in a single Get.
type Get_MultiResponse struct {
	Responses            []*Get_Response `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Get_MultiResponse) Reset()         { *m = Get_MultiResponse{} }
func (m *Get_MultiResponse) String() string { return proto.CompactTextString(m) }
func (*Get_MultiResponse) ProtoMessage()    {}
func (*Get_MultiResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{2, 3}
}
func (m *Get_MultiResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Get_MultiResponse.Unmarshal(m, b)
}
func (m *Get_MultiResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Get_MultiResponse.Marshal(b, m, deterministic)
}
func (dst *Get_MultiResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Get_MultiResponse.Merge(dst, src)
}
func (m *Get_MultiResponse) XXX_Size() int {
	return xxx_messageInfo_Get_MultiResponse.Size(m)
}
func (m *Get_MultiResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_Get_MultiResponse.DiscardUnknown(m)
}

var xxx_messageInfo_Get_MultiResponse proto.InternalMessageInfo

func (m *Get_MultiResponse) GetResponses() []*Get_Response {
	if m != nil {
		return m.Responses
	}
	return nil
}

// Close contains the structures for Close RPC calls.
type Close struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Close) Reset()         { *m = Close{} }
func (m *Close) String() string { return proto.CompactTextString(m) }
func (*Close) ProtoMessage()    {}
func (*Close) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{3}
}
func (m *Close) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Close.Unmarshal(m, b)
}
func (m *Close) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Close.Marshal(b, m, deterministic)
}
func (dst *Close) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Close.Merge(dst, src)
}
func (m *Close) XXX_Size() int {
	return xxx_messageInfo_Close.Size(m)
}
func (m *Close) XXX_DiscardUnknown() {
	xxx_messageInfo_Close.DiscardUnknown(m)
}

var xxx_messageInfo_Close proto.InternalMessageInfo

type Close_Request struct {
	InstanceId           uint64   `protobuf:"varint,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Close_Request) Reset()         { *m = Close_Request{} }
func (m *Close_Request) String() string { return proto.CompactTextString(m) }
func (*Close_Request) ProtoMessage()    {}
func (*Close_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{3, 0}
}
func (m *Close_Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Close_Request.Unmarshal(m, b)
}
func (m *Close_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Close_Request.Marshal(b, m, deterministic)
}
func (dst *Close_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Close_Request.Merge(dst, src)
}
func (m *Close_Request) XXX_Size() int {
	return xxx_messageInfo_Close_Request.Size(m)
}
func (m *Close_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Close_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Close_Request proto.InternalMessageInfo

func (m *Close_Request) GetInstanceId() uint64 {
	if m != nil {
		return m.InstanceId
	}
	return 0
}

// Value represents a Sentinel value.
type Value struct {
	// type is the type of this value
	Type Value_Type `protobuf:"varint,1,opt,name=type,proto3,enum=proto.Value_Type" json:"type,omitempty"`
	// value is the value only if the type is not UNDEFINED or NULL.
	// If the value is UNDEFINED or NULL, then the value is known.
	//
	// Types that are valid to be assigned to Value:
	//	*Value_ValueBool
	//	*Value_ValueInt
	//	*Value_ValueFloat
	//	*Value_ValueString
	//	*Value_ValueList
	//	*Value_ValueMap
	Value                isValue_Value `protobuf_oneof:"value"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{4}
}
func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (dst *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(dst, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetType() Value_Type {
	if m != nil {
		return m.Type
	}
	return Value_INVALID
}

type isValue_Value interface {
	isValue_Value()
}

type Value_ValueBool struct {
	ValueBool bool `protobuf:"varint,2,opt,name=value_bool,json=valueBool,proto3,oneof"`
}

type Value_ValueInt struct {
	ValueInt int64 `protobuf:"varint,3,opt,name=value_int,json=valueInt,proto3,oneof"`
}

type Value_ValueFloat struct {
	ValueFloat float64 `protobuf:"fixed64,4,opt,name=value_float,json=valueFloat,proto3,oneof"`
}

type Value_ValueString struct {
	ValueString string `protobuf:"bytes,5,opt,name=value_string,json=valueString,proto3,oneof"`
}

type Value_ValueList struct {
	ValueList *Value_List `protobuf:"bytes,6,opt,name=value_list,json=valueList,proto3,oneof"`
}

type Value_ValueMap struct {
	ValueMap *Value_Map `protobuf:"bytes,7,opt,name=value_map,json=valueMap,proto3,oneof"`
}

func (*Value_ValueBool) isValue_Value() {}

func (*Value_ValueInt) isValue_Value() {}

func (*Value_ValueFloat) isValue_Value() {}

func (*Value_ValueString) isValue_Value() {}

func (*Value_ValueList) isValue_Value() {}

func (*Value_ValueMap) isValue_Value() {}

func (m *Value) GetValue() isValue_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Value) GetValueBool() bool {
	if x, ok := m.GetValue().(*Value_ValueBool); ok {
		return x.ValueBool
	}
	return false
}

func (m *Value) GetValueInt() int64 {
	if x, ok := m.GetValue().(*Value_ValueInt); ok {
		return x.ValueInt
	}
	return 0
}

func (m *Value) GetValueFloat() float64 {
	if x, ok := m.GetValue().(*Value_ValueFloat); ok {
		return x.ValueFloat
	}
	return 0
}

func (m *Value) GetValueString() string {
	if x, ok := m.GetValue().(*Value_ValueString); ok {
		return x.ValueString
	}
	return ""
}

func (m *Value) GetValueList() *Value_List {
	if x, ok := m.GetValue().(*Value_ValueList); ok {
		return x.ValueList
	}
	return nil
}

func (m *Value) GetValueMap() *Value_Map {
	if x, ok := m.GetValue().(*Value_ValueMap); ok {
		return x.ValueMap
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Value) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Value_OneofMarshaler, _Value_OneofUnmarshaler, _Value_OneofSizer, []interface{}{
		(*Value_ValueBool)(nil),
		(*Value_ValueInt)(nil),
		(*Value_ValueFloat)(nil),
		(*Value_ValueString)(nil),
		(*Value_ValueList)(nil),
		(*Value_ValueMap)(nil),
	}
}

func _Value_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Value)
	// value
	switch x := m.Value.(type) {
	case *Value_ValueBool:
		t := uint64(0)
		if x.ValueBool {
			t = 1
		}
		b.EncodeVarint(2<<3 | proto.WireVarint)
		b.EncodeVarint(t)
	case *Value_ValueInt:
		b.EncodeVarint(3<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.ValueInt))
	case *Value_ValueFloat:
		b.EncodeVarint(4<<3 | proto.WireFixed64)
		b.EncodeFixed64(math.Float64bits(x.ValueFloat))
	case *Value_ValueString:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.ValueString)
	case *Value_ValueList:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ValueList); err != nil {
			return err
		}
	case *Value_ValueMap:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ValueMap); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Value.Value has unexpected type %T", x)
	}
	return nil
}

func _Value_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Value)
	switch tag {
	case 2: // value.value_bool
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Value = &Value_ValueBool{x != 0}
		return true, err
	case 3: // value.value_int
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Value = &Value_ValueInt{int64(x)}
		return true, err
	case 4: // value.value_float
		if wire != proto.WireFixed64 {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeFixed64()
		m.Value = &Value_ValueFloat{math.Float64frombits(x)}
		return true, err
	case 5: // value.value_string
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Value = &Value_ValueString{x}
		return true, err
	case 6: // value.value_list
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Value_List)
		err := b.DecodeMessage(msg)
		m.Value = &Value_ValueList{msg}
		return true, err
	case 7: // value.value_map
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Value_Map)
		err := b.DecodeMessage(msg)
		m.Value = &Value_ValueMap{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Value_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Value)
	// value
	switch x := m.Value.(type) {
	case *Value_ValueBool:
		n += 1 // tag and wire
		n += 1
	case *Value_ValueInt:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(x.ValueInt))
	case *Value_ValueFloat:
		n += 1 // tag and wire
		n += 8
	case *Value_ValueString:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.ValueString)))
		n += len(x.ValueString)
	case *Value_ValueList:
		s := proto.Size(x.ValueList)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Value_ValueMap:
		s := proto.Size(x.ValueMap)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Value_KV struct {
	Key                  *Value   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                *Value   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Value_KV) Reset()         { *m = Value_KV{} }
func (m *Value_KV) String() string { return proto.CompactTextString(m) }
func (*Value_KV) ProtoMessage()    {}
func (*Value_KV) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{4, 0}
}
func (m *Value_KV) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value_KV.Unmarshal(m, b)
}
func (m *Value_KV) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value_KV.Marshal(b, m, deterministic)
}
func (dst *Value_KV) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value_KV.Merge(dst, src)
}
func (m *Value_KV) XXX_Size() int {
	return xxx_messageInfo_Value_KV.Size(m)
}
func (m *Value_KV) XXX_DiscardUnknown() {
	xxx_messageInfo_Value_KV.DiscardUnknown(m)
}

var xxx_messageInfo_Value_KV proto.InternalMessageInfo

func (m *Value_KV) GetKey() *Value {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Value_KV) GetValue() *Value {
	if m != nil {
		return m.Value
	}
	return nil
}

type Value_Map struct {
	Elems                []*Value_KV `protobuf:"bytes,1,rep,name=elems,proto3" json:"elems,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Value_Map) Reset()         { *m = Value_Map{} }
func (m *Value_Map) String() string { return proto.CompactTextString(m) }
func (*Value_Map) ProtoMessage()    {}
func (*Value_Map) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{4, 1}
}
func (m *Value_Map) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value_Map.Unmarshal(m, b)
}
func (m *Value_Map) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value_Map.Marshal(b, m, deterministic)
}
func (dst *Value_Map) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value_Map.Merge(dst, src)
}
func (m *Value_Map) XXX_Size() int {
	return xxx_messageInfo_Value_Map.Size(m)
}
func (m *Value_Map) XXX_DiscardUnknown() {
	xxx_messageInfo_Value_Map.DiscardUnknown(m)
}

var xxx_messageInfo_Value_Map proto.InternalMessageInfo

func (m *Value_Map) GetElems() []*Value_KV {
	if m != nil {
		return m.Elems
	}
	return nil
}

type Value_List struct {
	Elems                []*Value `protobuf:"bytes,1,rep,name=elems,proto3" json:"elems,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Value_List) Reset()         { *m = Value_List{} }
func (m *Value_List) String() string { return proto.CompactTextString(m) }
func (*Value_List) ProtoMessage()    {}
func (*Value_List) Descriptor() ([]byte, []int) {
	return fileDescriptor_import_a41f914c4f8436db, []int{4, 2}
}
func (m *Value_List) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value_List.Unmarshal(m, b)
}
func (m *Value_List) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value_List.Marshal(b, m, deterministic)
}
func (dst *Value_List) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value_List.Merge(dst, src)
}
func (m *Value_List) XXX_Size() int {
	return xxx_messageInfo_Value_List.Size(m)
}
func (m *Value_List) XXX_DiscardUnknown() {
	xxx_messageInfo_Value_List.DiscardUnknown(m)
}

var xxx_messageInfo_Value_List proto.InternalMessageInfo

func (m *Value_List) GetElems() []*Value {
	if m != nil {
		return m.Elems
	}
	return nil
}

func init() {
	proto.RegisterType((*EmptyResp)(nil), "proto.EmptyResp")
	proto.RegisterType((*Configure)(nil), "proto.Configure")
	proto.RegisterType((*Configure_Request)(nil), "proto.Configure.Request")
	proto.RegisterType((*Configure_Response)(nil), "proto.Configure.Response")
	proto.RegisterType((*Get)(nil), "proto.Get")
	proto.RegisterType((*Get_Request)(nil), "proto.Get.Request")
	proto.RegisterType((*Get_Response)(nil), "proto.Get.Response")
	proto.RegisterType((*Get_MultiRequest)(nil), "proto.Get.MultiRequest")
	proto.RegisterType((*Get_MultiResponse)(nil), "proto.Get.MultiResponse")
	proto.RegisterType((*Close)(nil), "proto.Close")
	proto.RegisterType((*Close_Request)(nil), "proto.Close.Request")
	proto.RegisterType((*Value)(nil), "proto.Value")
	proto.RegisterType((*Value_KV)(nil), "proto.Value.KV")
	proto.RegisterType((*Value_Map)(nil), "proto.Value.Map")
	proto.RegisterType((*Value_List)(nil), "proto.Value.List")
	proto.RegisterEnum("proto.Value_Type", Value_Type_name, Value_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ImportClient is the client API for Import service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ImportClient interface {
	Configure(ctx context.Context, in *Configure_Request, opts ...grpc.CallOption) (*Configure_Response, error)
	Get(ctx context.Context, in *Get_MultiRequest, opts ...grpc.CallOption) (*Get_MultiResponse, error)
	Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*EmptyResp, error)
}

type importClient struct {
	cc *grpc.ClientConn
}

func NewImportClient(cc *grpc.ClientConn) ImportClient {
	return &importClient{cc}
}

func (c *importClient) Configure(ctx context.Context, in *Configure_Request, opts ...grpc.CallOption) (*Configure_Response, error) {
	out := new(Configure_Response)
	err := c.cc.Invoke(ctx, "/proto.Import/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *importClient) Get(ctx context.Context, in *Get_MultiRequest, opts ...grpc.CallOption) (*Get_MultiResponse, error) {
	out := new(Get_MultiResponse)
	err := c.cc.Invoke(ctx, "/proto.Import/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *importClient) Close(ctx context.Context, in *Close_Request, opts ...grpc.CallOption) (*EmptyResp, error) {
	out := new(EmptyResp)
	err := c.cc.Invoke(ctx, "/proto.Import/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImportServer is the server API for Import service.
type ImportServer interface {
	Configure(context.Context, *Configure_Request) (*Configure_Response, error)
	Get(context.Context, *Get_MultiRequest) (*Get_MultiResponse, error)
	Close(context.Context, *Close_Request) (*EmptyResp, error)
}

func RegisterImportServer(s *grpc.Server, srv ImportServer) {
	s.RegisterService(&_Import_serviceDesc, srv)
}

func _Import_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Configure_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImportServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Import/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImportServer).Configure(ctx, req.(*Configure_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Import_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Get_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImportServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Import/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImportServer).Get(ctx, req.(*Get_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Import_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Close_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImportServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Import/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImportServer).Close(ctx, req.(*Close_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Import_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Import",
	HandlerType: (*ImportServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Configure",
			Handler:    _Import_Configure_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Import_Get_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Import_Close_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "import.proto",
}

func init() { proto.RegisterFile("import.proto", fileDescriptor_import_a41f914c4f8436db) }

var fileDescriptor_import_a41f914c4f8436db = []byte{
	// 694 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdd, 0x6e, 0xda, 0x48,
	0x14, 0xc6, 0xf8, 0x0f, 0x1f, 0xc8, 0xae, 0x77, 0x76, 0x57, 0x71, 0x2d, 0xb5, 0xa1, 0x4e, 0x23,
	0xa1, 0xb4, 0x22, 0x2a, 0xb9, 0xe9, 0x55, 0xd5, 0x10, 0x92, 0x60, 0x05, 0x48, 0x35, 0x21, 0xdc,
	0x46, 0x0e, 0x4c, 0x22, 0x0b, 0x63, 0xbb, 0x78, 0xa8, 0xea, 0xbe, 0x56, 0xd5, 0x87, 0xe8, 0x65,
	0xdf, 0xa8, 0x9a, 0x63, 0x9b, 0x00, 0x49, 0xd5, 0xf6, 0xca, 0x33, 0xdf, 0x77, 0x7e, 0xbe, 0xf3,
	0xe3, 0x81, 0x9a, 0x3f, 0x8b, 0xa3, 0x39, 0x6f, 0xc6, 0xf3, 0x88, 0x47, 0x44, 0xc5, 0x8f, 0x53,
	0x05, 0xe3, 0x64, 0x16, 0xf3, 0x94, 0xb2, 0x24, 0x76, 0x7c, 0x30, 0x8e, 0xa3, 0xf0, 0xd6, 0xbf,
	0x5b, 0xcc, 0x99, 0x7d, 0x00, 0x3a, 0x65, 0x1f, 0x16, 0x2c, 0xe1, 0xe4, 0x05, 0x68, 0x63, 0xc4,
	0x2d, 0xb9, 0x2e, 0x35, 0xaa, 0xad, 0x5a, 0x16, 0xa3, 0x39, 0xf2, 0x82, 0x05, 0xa3, 0x39, 0x67,
	0xbf, 0x84, 0x8a, 0x88, 0x12, 0x85, 0x09, 0x23, 0x3b, 0x50, 0xf5, 0xc3, 0x84, 0x7b, 0xe1, 0x98,
	0x5d, 0xfb, 0x13, 0x4b, 0xaa, 0x4b, 0x0d, 0x85, 0x42, 0x01, 0xb9, 0x13, 0xe7, 0xbb, 0x0c, 0xf2,
	0x19, 0xe3, 0xf6, 0x37, 0xe9, 0x3e, 0xcd, 0xaf, 0x9c, 0xc8, 0x36, 0xe8, 0xec, 0x13, 0x1b, 0x0b,
	0xb2, 0x8c, 0xa4, 0x26, 0xae, 0xee, 0x84, 0xec, 0xc2, 0x16, 0x12, 0x13, 0xe6, 0x4d, 0x02, 0x3f,
	0x64, 0xa8, 0x53, 0xa1, 0x35, 0x01, 0x76, 0x72, 0x8c, 0x10, 0x50, 0xa6, 0x2c, 0x4d, 0x2c, 0xa5,
	0x2e, 0x37, 0x0c, 0x8a, 0x67, 0xf2, 0x3f, 0x68, 0x53, 0x96, 0x8a, 0x80, 0x2a, 0x7a, 0xa8, 0x53,
	0x96, 0xba, 0x13, 0x61, 0x3a, 0xf6, 0x82, 0xc0, 0xd2, 0xea, 0x52, 0xa3, 0x42, 0xf1, 0x4c, 0xea,
	0xa0, 0x78, 0xf3, 0xbb, 0xc4, 0xd2, 0xeb, 0xf2, 0x83, 0x16, 0x20, 0x63, 0x7f, 0xfe, 0x83, 0x06,
	0xac, 0x64, 0x2e, 0x6f, 0x64, 0x46, 0x91, 0xf2, 0x8a, 0x48, 0x07, 0xd4, 0x8f, 0x22, 0x8d, 0xa5,
	0x3c, 0xd2, 0xfd, 0x8c, 0xb2, 0xdf, 0x42, 0xad, 0xbf, 0x08, 0xb8, 0x5f, 0xf4, 0xb2, 0x09, 0x95,
	0x79, 0x76, 0x4c, 0x2c, 0x09, 0x15, 0x93, 0xdc, 0xed, 0x8c, 0xf1, 0x66, 0x6e, 0x45, 0x97, 0x36,
	0x76, 0x1b, 0xb6, 0x72, 0xff, 0xbc, 0x80, 0xd7, 0x60, 0xcc, 0xf3, 0x73, 0x11, 0xe1, 0xdf, 0xb5,
	0x08, 0x19, 0x47, 0xef, 0xad, 0x9c, 0x43, 0x50, 0x8f, 0x83, 0x28, 0x61, 0xf6, 0xfe, 0xef, 0xcf,
	0xd4, 0xf9, 0xa2, 0x80, 0x8a, 0x95, 0x90, 0x3d, 0x50, 0x78, 0x1a, 0x33, 0xb4, 0xf9, 0xab, 0xf5,
	0xcf, 0x6a, 0x95, 0xcd, 0x61, 0x1a, 0x33, 0x8a, 0x34, 0xd9, 0x01, 0xc0, 0x92, 0xaf, 0x6f, 0xa2,
	0x28, 0xc0, 0xe6, 0x55, 0xba, 0x25, 0x6a, 0x20, 0xd6, 0x8e, 0xa2, 0x80, 0x3c, 0x85, 0xec, 0x72,
	0xed, 0x87, 0x1c, 0x17, 0x41, 0xee, 0x96, 0x68, 0x05, 0x21, 0x37, 0xe4, 0xe4, 0x39, 0x54, 0x33,
	0xfa, 0x36, 0x88, 0x3c, 0x8e, 0x3d, 0x95, 0xba, 0x25, 0x9a, 0x05, 0x3d, 0x15, 0x18, 0xd9, 0x85,
	0x5a, 0x66, 0x92, 0xf0, 0xb9, 0x1f, 0xde, 0xe1, 0x6e, 0x18, 0xdd, 0x12, 0xcd, 0x1c, 0x2f, 0x11,
	0x24, 0xad, 0x42, 0x47, 0xe0, 0x27, 0x1c, 0x37, 0xa5, 0xba, 0x21, 0xba, 0xe7, 0x27, 0x7c, 0x29,
	0x4d, 0x5c, 0xc8, 0x41, 0x21, 0x6d, 0xe6, 0xc5, 0x96, 0x8e, 0x2e, 0xe6, 0x9a, 0x4b, 0xdf, 0x8b,
	0x97, 0x62, 0xfb, 0x5e, 0x6c, 0x77, 0xa1, 0x7c, 0x3e, 0x22, 0xcf, 0x40, 0x9e, 0xb2, 0x14, 0x1b,
	0xb3, 0x39, 0x7e, 0x41, 0xdc, 0x2f, 0x48, 0xf9, 0xe7, 0x0b, 0xf2, 0x0a, 0xe4, 0xbe, 0x17, 0x93,
	0x3d, 0x50, 0x59, 0xc0, 0x66, 0xc5, 0x48, 0xff, 0x5e, 0xcb, 0x7e, 0x3e, 0xa2, 0x19, 0x6b, 0xef,
	0x83, 0x82, 0x82, 0x9d, 0x75, 0xf3, 0x8d, 0xc8, 0x48, 0x39, 0x3e, 0x28, 0x62, 0x3c, 0xa4, 0x0a,
	0xba, 0x3b, 0x18, 0x1d, 0xf5, 0xdc, 0x8e, 0x59, 0x22, 0x5b, 0x60, 0x5c, 0x0d, 0x3a, 0x27, 0xa7,
	0xee, 0xe0, 0xa4, 0x63, 0x4a, 0xa4, 0x02, 0xca, 0xe0, 0xaa, 0xd7, 0x33, 0xcb, 0xe2, 0xd4, 0xbe,
	0xb8, 0xe8, 0x99, 0x32, 0xd1, 0x41, 0x76, 0x07, 0x43, 0x53, 0x21, 0x06, 0xa8, 0xa7, 0xbd, 0x8b,
	0xa3, 0xa1, 0xa9, 0x12, 0x00, 0xed, 0x72, 0x48, 0xdd, 0xc1, 0x99, 0xa9, 0x09, 0xcb, 0x9e, 0x7b,
	0x39, 0x34, 0x75, 0x61, 0xd9, 0x3f, 0x7a, 0x6f, 0x56, 0xda, 0x7a, 0x5e, 0x68, 0xeb, 0xab, 0x04,
	0x9a, 0x8b, 0xcf, 0x19, 0x79, 0xb7, 0xf2, 0x68, 0x11, 0x2b, 0x17, 0xb8, 0x44, 0x8a, 0x55, 0xb7,
	0x9f, 0x3c, 0xc2, 0xe4, 0xab, 0xfe, 0x06, 0x9f, 0x22, 0xb2, 0xbd, 0xb2, 0xde, 0xab, 0xff, 0x92,
	0x6d, 0x3d, 0x24, 0x72, 0xcf, 0x83, 0x7c, 0xe3, 0xc9, 0x7f, 0x45, 0x74, 0x71, 0x5b, 0xe6, 0x2c,
	0x66, 0xbb, 0x7c, 0x61, 0x6f, 0x34, 0x04, 0x0e, 0x7f, 0x04, 0x00, 0x00, 0xff, 0xff, 0x1f, 0xff,
	0x61, 0xb8, 0x8c, 0x05, 0x00, 0x00,
}
