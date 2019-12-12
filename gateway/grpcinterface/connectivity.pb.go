// Code generated by protoc-gen-go. DO NOT EDIT.
// source: connectivity.proto

package trueno

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PingReply struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingReply) Reset()         { *m = PingReply{} }
func (m *PingReply) String() string { return proto.CompactTextString(m) }
func (*PingReply) ProtoMessage()    {}
func (*PingReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_2872c2021a21e8fe, []int{0}
}

func (m *PingReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingReply.Unmarshal(m, b)
}
func (m *PingReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingReply.Marshal(b, m, deterministic)
}
func (m *PingReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingReply.Merge(m, src)
}
func (m *PingReply) XXX_Size() int {
	return xxx_messageInfo_PingReply.Size(m)
}
func (m *PingReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PingReply.DiscardUnknown(m)
}

var xxx_messageInfo_PingReply proto.InternalMessageInfo

func (m *PingReply) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type ResourcesReply struct {
	Cpu                  string   `protobuf:"bytes,1,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Mem                  string   `protobuf:"bytes,2,opt,name=mem,proto3" json:"mem,omitempty"`
	Gpu                  string   `protobuf:"bytes,3,opt,name=gpu,proto3" json:"gpu,omitempty"`
	Dsk                  string   `protobuf:"bytes,4,opt,name=dsk,proto3" json:"dsk,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResourcesReply) Reset()         { *m = ResourcesReply{} }
func (m *ResourcesReply) String() string { return proto.CompactTextString(m) }
func (*ResourcesReply) ProtoMessage()    {}
func (*ResourcesReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_2872c2021a21e8fe, []int{1}
}

func (m *ResourcesReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourcesReply.Unmarshal(m, b)
}
func (m *ResourcesReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourcesReply.Marshal(b, m, deterministic)
}
func (m *ResourcesReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourcesReply.Merge(m, src)
}
func (m *ResourcesReply) XXX_Size() int {
	return xxx_messageInfo_ResourcesReply.Size(m)
}
func (m *ResourcesReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourcesReply.DiscardUnknown(m)
}

var xxx_messageInfo_ResourcesReply proto.InternalMessageInfo

func (m *ResourcesReply) GetCpu() string {
	if m != nil {
		return m.Cpu
	}
	return ""
}

func (m *ResourcesReply) GetMem() string {
	if m != nil {
		return m.Mem
	}
	return ""
}

func (m *ResourcesReply) GetGpu() string {
	if m != nil {
		return m.Gpu
	}
	return ""
}

func (m *ResourcesReply) GetDsk() string {
	if m != nil {
		return m.Dsk
	}
	return ""
}

func init() {
	proto.RegisterType((*PingReply)(nil), "trueno.PingReply")
	proto.RegisterType((*ResourcesReply)(nil), "trueno.ResourcesReply")
}

func init() { proto.RegisterFile("connectivity.proto", fileDescriptor_2872c2021a21e8fe) }

var fileDescriptor_2872c2021a21e8fe = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0xd0, 0xcf, 0x4a, 0xc4, 0x30,
	0x10, 0xc7, 0xf1, 0xad, 0xbb, 0xac, 0xec, 0xb0, 0x88, 0x1b, 0x41, 0xc2, 0x9e, 0x24, 0x20, 0x78,
	0x2a, 0xa2, 0x6f, 0xa0, 0x57, 0x11, 0xe9, 0xc9, 0xb3, 0xe9, 0x50, 0x42, 0x9b, 0x4c, 0xcc, 0x9f,
	0x42, 0xef, 0x3e, 0xb8, 0xc4, 0xb1, 0x62, 0x61, 0x6f, 0xed, 0x87, 0x5f, 0xe0, 0x9b, 0x80, 0xd0,
	0xe4, 0x1c, 0xea, 0x64, 0x46, 0x93, 0xa6, 0xda, 0x07, 0x4a, 0x24, 0xb6, 0x29, 0x64, 0x74, 0x74,
	0xdc, 0x6b, 0xb2, 0x96, 0x1c, 0xab, 0xba, 0x85, 0xdd, 0x9b, 0x71, 0x5d, 0x83, 0x7e, 0x98, 0x84,
	0x84, 0xf3, 0x11, 0x43, 0x34, 0xe4, 0x64, 0x75, 0x53, 0xdd, 0xed, 0x9a, 0xf9, 0x57, 0xbd, 0xc3,
	0x45, 0x83, 0x91, 0x72, 0xd0, 0x18, 0x79, 0x7b, 0x09, 0x6b, 0xed, 0xf3, 0xef, 0xae, 0x7c, 0x16,
	0xb1, 0x68, 0xe5, 0x19, 0x8b, 0x45, 0x5b, 0xa4, 0xf3, 0x59, 0xae, 0x59, 0x3a, 0xde, 0xb4, 0xb1,
	0x97, 0x1b, 0x96, 0x36, 0xf6, 0x0f, 0x5f, 0x15, 0xec, 0x9f, 0xff, 0xd5, 0x8a, 0x7b, 0xd8, 0x94,
	0x22, 0x71, 0x55, 0x73, 0x70, 0xcd, 0x7d, 0x9f, 0x19, 0x63, 0x3a, 0x1e, 0x96, 0xe8, 0x87, 0x49,
	0xad, 0xc4, 0x13, 0x1c, 0x5e, 0x4c, 0x4c, 0xaf, 0xd4, 0xe2, 0x5f, 0xe4, 0xe9, 0xe3, 0xd7, 0x33,
	0x2e, 0x2f, 0xa3, 0x56, 0x1f, 0xdb, 0x9f, 0xe7, 0x78, 0xfc, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xdc,
	0xe5, 0xa6, 0x52, 0x3a, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConnectivityClient is the client API for Connectivity service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConnectivityClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error)
	ListNodeResources(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*ResourcesReply, error)
}

type connectivityClient struct {
	cc *grpc.ClientConn
}

func NewConnectivityClient(cc *grpc.ClientConn) ConnectivityClient {
	return &connectivityClient{cc}
}

func (c *connectivityClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReply, error) {
	out := new(PingReply)
	err := c.cc.Invoke(ctx, "/trueno.Connectivity/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectivityClient) ListNodeResources(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*ResourcesReply, error) {
	out := new(ResourcesReply)
	err := c.cc.Invoke(ctx, "/trueno.Connectivity/ListNodeResources", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectivityServer is the server API for Connectivity service.
type ConnectivityServer interface {
	Ping(context.Context, *PingRequest) (*PingReply, error)
	ListNodeResources(context.Context, *PingRequest) (*ResourcesReply, error)
}

// UnimplementedConnectivityServer can be embedded to have forward compatible implementations.
type UnimplementedConnectivityServer struct {
}

func (*UnimplementedConnectivityServer) Ping(ctx context.Context, req *PingRequest) (*PingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedConnectivityServer) ListNodeResources(ctx context.Context, req *PingRequest) (*ResourcesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNodeResources not implemented")
}

func RegisterConnectivityServer(s *grpc.Server, srv ConnectivityServer) {
	s.RegisterService(&_Connectivity_serviceDesc, srv)
}

func _Connectivity_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectivityServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trueno.Connectivity/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectivityServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Connectivity_ListNodeResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectivityServer).ListNodeResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trueno.Connectivity/ListNodeResources",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectivityServer).ListNodeResources(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Connectivity_serviceDesc = grpc.ServiceDesc{
	ServiceName: "trueno.Connectivity",
	HandlerType: (*ConnectivityServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Connectivity_Ping_Handler,
		},
		{
			MethodName: "ListNodeResources",
			Handler:    _Connectivity_ListNodeResources_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "connectivity.proto",
}
