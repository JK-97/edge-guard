// Code generated by protoc-gen-go. DO NOT EDIT.
// source: backend.proto

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

type SupportedReply struct {
	Support              string   `protobuf:"bytes,1,opt,name=support,proto3" json:"support,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SupportedReply) Reset()         { *m = SupportedReply{} }
func (m *SupportedReply) String() string { return proto.CompactTextString(m) }
func (*SupportedReply) ProtoMessage()    {}
func (*SupportedReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ab9ba5b8d8b2ba5, []int{0}
}

func (m *SupportedReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SupportedReply.Unmarshal(m, b)
}
func (m *SupportedReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SupportedReply.Marshal(b, m, deterministic)
}
func (m *SupportedReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SupportedReply.Merge(m, src)
}
func (m *SupportedReply) XXX_Size() int {
	return xxx_messageInfo_SupportedReply.Size(m)
}
func (m *SupportedReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SupportedReply.DiscardUnknown(m)
}

var xxx_messageInfo_SupportedReply proto.InternalMessageInfo

func (m *SupportedReply) GetSupport() string {
	if m != nil {
		return m.Support
	}
	return ""
}

type RunningReply struct {
	Status               []*RunningReply_Status `protobuf:"bytes,1,rep,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *RunningReply) Reset()         { *m = RunningReply{} }
func (m *RunningReply) String() string { return proto.CompactTextString(m) }
func (*RunningReply) ProtoMessage()    {}
func (*RunningReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ab9ba5b8d8b2ba5, []int{1}
}

func (m *RunningReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunningReply.Unmarshal(m, b)
}
func (m *RunningReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunningReply.Marshal(b, m, deterministic)
}
func (m *RunningReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunningReply.Merge(m, src)
}
func (m *RunningReply) XXX_Size() int {
	return xxx_messageInfo_RunningReply.Size(m)
}
func (m *RunningReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RunningReply.DiscardUnknown(m)
}

var xxx_messageInfo_RunningReply proto.InternalMessageInfo

func (m *RunningReply) GetStatus() []*RunningReply_Status {
	if m != nil {
		return m.Status
	}
	return nil
}

type RunningReply_Status struct {
	Bid                  string   `protobuf:"bytes,1,opt,name=bid,proto3" json:"bid,omitempty"`
	Model                string   `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty"`
	Status               string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Msg                  string   `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunningReply_Status) Reset()         { *m = RunningReply_Status{} }
func (m *RunningReply_Status) String() string { return proto.CompactTextString(m) }
func (*RunningReply_Status) ProtoMessage()    {}
func (*RunningReply_Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ab9ba5b8d8b2ba5, []int{1, 0}
}

func (m *RunningReply_Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunningReply_Status.Unmarshal(m, b)
}
func (m *RunningReply_Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunningReply_Status.Marshal(b, m, deterministic)
}
func (m *RunningReply_Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunningReply_Status.Merge(m, src)
}
func (m *RunningReply_Status) XXX_Size() int {
	return xxx_messageInfo_RunningReply_Status.Size(m)
}
func (m *RunningReply_Status) XXX_DiscardUnknown() {
	xxx_messageInfo_RunningReply_Status.DiscardUnknown(m)
}

var xxx_messageInfo_RunningReply_Status proto.InternalMessageInfo

func (m *RunningReply_Status) GetBid() string {
	if m != nil {
		return m.Bid
	}
	return ""
}

func (m *RunningReply_Status) GetModel() string {
	if m != nil {
		return m.Model
	}
	return ""
}

func (m *RunningReply_Status) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *RunningReply_Status) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type BackendRequest struct {
	Bid                  string   `protobuf:"bytes,1,opt,name=bid,proto3" json:"bid,omitempty"`
	Btype                string   `protobuf:"bytes,2,opt,name=btype,proto3" json:"btype,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BackendRequest) Reset()         { *m = BackendRequest{} }
func (m *BackendRequest) String() string { return proto.CompactTextString(m) }
func (*BackendRequest) ProtoMessage()    {}
func (*BackendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ab9ba5b8d8b2ba5, []int{2}
}

func (m *BackendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackendRequest.Unmarshal(m, b)
}
func (m *BackendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackendRequest.Marshal(b, m, deterministic)
}
func (m *BackendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackendRequest.Merge(m, src)
}
func (m *BackendRequest) XXX_Size() int {
	return xxx_messageInfo_BackendRequest.Size(m)
}
func (m *BackendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BackendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BackendRequest proto.InternalMessageInfo

func (m *BackendRequest) GetBid() string {
	if m != nil {
		return m.Bid
	}
	return ""
}

func (m *BackendRequest) GetBtype() string {
	if m != nil {
		return m.Btype
	}
	return ""
}

func init() {
	proto.RegisterType((*SupportedReply)(nil), "trueno.SupportedReply")
	proto.RegisterType((*RunningReply)(nil), "trueno.RunningReply")
	proto.RegisterType((*RunningReply_Status)(nil), "trueno.RunningReply.Status")
	proto.RegisterType((*BackendRequest)(nil), "trueno.BackendRequest")
}

func init() { proto.RegisterFile("backend.proto", fileDescriptor_5ab9ba5b8d8b2ba5) }

var fileDescriptor_5ab9ba5b8d8b2ba5 = []byte{
	// 322 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0xbd, 0x6e, 0xf2, 0x40,
	0x10, 0xc4, 0xf0, 0x7d, 0x46, 0x59, 0x08, 0x82, 0x03, 0x21, 0xcb, 0x69, 0x90, 0x2b, 0x94, 0xc2,
	0x05, 0x34, 0x69, 0x83, 0xd2, 0x44, 0x4a, 0x11, 0x19, 0xaa, 0x74, 0x36, 0x5e, 0xa1, 0x53, 0xec,
	0x3b, 0xc7, 0x77, 0x57, 0x90, 0x97, 0xc9, 0x9b, 0xe5, 0x59, 0x22, 0xdf, 0x8f, 0x15, 0x24, 0x2b,
	0x91, 0xd2, 0x79, 0x47, 0x3b, 0x33, 0x3b, 0xe3, 0x83, 0xeb, 0x2c, 0x3d, 0xbe, 0x22, 0xcb, 0xe3,
	0xaa, 0xe6, 0x92, 0x13, 0x5f, 0xd6, 0x0a, 0x19, 0x0f, 0xc7, 0x47, 0x5e, 0x96, 0x9c, 0x19, 0x34,
	0xba, 0x85, 0xc9, 0x5e, 0x55, 0x15, 0xaf, 0x25, 0xe6, 0x09, 0x56, 0xc5, 0x99, 0x04, 0x30, 0x14,
	0x06, 0x09, 0xbc, 0x95, 0xb7, 0xbe, 0x4a, 0xdc, 0x18, 0x7d, 0x78, 0x30, 0x4e, 0x14, 0x63, 0x94,
	0x9d, 0xcc, 0xea, 0x16, 0x7c, 0x21, 0x53, 0xa9, 0x44, 0xe0, 0xad, 0x06, 0xeb, 0xd1, 0xe6, 0x26,
	0x36, 0x1e, 0xf1, 0xf7, 0xad, 0x78, 0xaf, 0x57, 0x12, 0xbb, 0x1a, 0xbe, 0x80, 0x6f, 0x10, 0x32,
	0x85, 0x41, 0x46, 0x73, 0xeb, 0xd2, 0x7c, 0x92, 0x05, 0xfc, 0x2f, 0x79, 0x8e, 0x45, 0xd0, 0xd7,
	0x98, 0x19, 0xc8, 0xb2, 0xb5, 0x19, 0x68, 0xd8, 0x4e, 0x0d, 0xbf, 0x14, 0xa7, 0xe0, 0x9f, 0xe1,
	0x97, 0xe2, 0x14, 0xdd, 0xc1, 0x64, 0x67, 0x42, 0x27, 0xf8, 0xa6, 0x50, 0xc8, 0x6e, 0x8f, 0x4c,
	0x9e, 0x2b, 0x74, 0x1e, 0x7a, 0xd8, 0x7c, 0xf6, 0x61, 0x68, 0xa9, 0x64, 0x07, 0xb3, 0x27, 0x2a,
	0x64, 0xdb, 0xcb, 0xe1, 0x5c, 0x21, 0x99, 0xbb, 0x6c, 0xcf, 0x3a, 0x98, 0x56, 0x0f, 0x97, 0x0e,
	0xbc, 0xec, 0x30, 0xea, 0x91, 0x1d, 0xcc, 0x1b, 0x0d, 0x5b, 0x84, 0x55, 0x16, 0xdd, 0x2a, 0x8b,
	0xae, 0xda, 0xb4, 0xc6, 0xec, 0x91, 0x51, 0x49, 0xd3, 0x82, 0xbe, 0xa3, 0x3b, 0xae, 0xb5, 0xbc,
	0x0c, 0x1a, 0xb6, 0xca, 0x09, 0x0a, 0x55, 0x48, 0xa7, 0xf1, 0x00, 0xa3, 0xe6, 0x8e, 0xdf, 0xd8,
	0x3f, 0xfd, 0xb9, 0xa8, 0x47, 0xee, 0x61, 0x7a, 0xc0, 0xba, 0xa4, 0x2c, 0x95, 0x7f, 0x3c, 0x24,
	0xf3, 0xf5, 0x7b, 0xdb, 0x7e, 0x05, 0x00, 0x00, 0xff, 0xff, 0xab, 0xd8, 0x7f, 0xbc, 0x96, 0x02,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BackendClient is the client API for Backend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BackendClient interface {
	ListSupportedType(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*SupportedReply, error)
	ListRunningBackends(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*RunningReply, error)
	InitializeBackend(ctx context.Context, in *BackendRequest, opts ...grpc.CallOption) (*ResultReply, error)
	ListBackend(ctx context.Context, in *BackendRequest, opts ...grpc.CallOption) (*RunningReply_Status, error)
	TerminateBackend(ctx context.Context, in *BackendRequest, opts ...grpc.CallOption) (*ResultReply, error)
}

type backendClient struct {
	cc *grpc.ClientConn
}

func NewBackendClient(cc *grpc.ClientConn) BackendClient {
	return &backendClient{cc}
}

func (c *backendClient) ListSupportedType(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*SupportedReply, error) {
	out := new(SupportedReply)
	err := c.cc.Invoke(ctx, "/trueno.Backend/ListSupportedType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) ListRunningBackends(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*RunningReply, error) {
	out := new(RunningReply)
	err := c.cc.Invoke(ctx, "/trueno.Backend/ListRunningBackends", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) InitializeBackend(ctx context.Context, in *BackendRequest, opts ...grpc.CallOption) (*ResultReply, error) {
	out := new(ResultReply)
	err := c.cc.Invoke(ctx, "/trueno.Backend/InitializeBackend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) ListBackend(ctx context.Context, in *BackendRequest, opts ...grpc.CallOption) (*RunningReply_Status, error) {
	out := new(RunningReply_Status)
	err := c.cc.Invoke(ctx, "/trueno.Backend/ListBackend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) TerminateBackend(ctx context.Context, in *BackendRequest, opts ...grpc.CallOption) (*ResultReply, error) {
	out := new(ResultReply)
	err := c.cc.Invoke(ctx, "/trueno.Backend/TerminateBackend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackendServer is the server API for Backend service.
type BackendServer interface {
	ListSupportedType(context.Context, *PingRequest) (*SupportedReply, error)
	ListRunningBackends(context.Context, *PingRequest) (*RunningReply, error)
	InitializeBackend(context.Context, *BackendRequest) (*ResultReply, error)
	ListBackend(context.Context, *BackendRequest) (*RunningReply_Status, error)
	TerminateBackend(context.Context, *BackendRequest) (*ResultReply, error)
}

// UnimplementedBackendServer can be embedded to have forward compatible implementations.
type UnimplementedBackendServer struct {
}

func (*UnimplementedBackendServer) ListSupportedType(ctx context.Context, req *PingRequest) (*SupportedReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSupportedType not implemented")
}
func (*UnimplementedBackendServer) ListRunningBackends(ctx context.Context, req *PingRequest) (*RunningReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRunningBackends not implemented")
}
func (*UnimplementedBackendServer) InitializeBackend(ctx context.Context, req *BackendRequest) (*ResultReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitializeBackend not implemented")
}
func (*UnimplementedBackendServer) ListBackend(ctx context.Context, req *BackendRequest) (*RunningReply_Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBackend not implemented")
}
func (*UnimplementedBackendServer) TerminateBackend(ctx context.Context, req *BackendRequest) (*ResultReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TerminateBackend not implemented")
}

func RegisterBackendServer(s *grpc.Server, srv BackendServer) {
	s.RegisterService(&_Backend_serviceDesc, srv)
}

func _Backend_ListSupportedType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).ListSupportedType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trueno.Backend/ListSupportedType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).ListSupportedType(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_ListRunningBackends_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).ListRunningBackends(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trueno.Backend/ListRunningBackends",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).ListRunningBackends(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_InitializeBackend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).InitializeBackend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trueno.Backend/InitializeBackend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).InitializeBackend(ctx, req.(*BackendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_ListBackend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).ListBackend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trueno.Backend/ListBackend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).ListBackend(ctx, req.(*BackendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_TerminateBackend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).TerminateBackend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/trueno.Backend/TerminateBackend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).TerminateBackend(ctx, req.(*BackendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backend_serviceDesc = grpc.ServiceDesc{
	ServiceName: "trueno.Backend",
	HandlerType: (*BackendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSupportedType",
			Handler:    _Backend_ListSupportedType_Handler,
		},
		{
			MethodName: "ListRunningBackends",
			Handler:    _Backend_ListRunningBackends_Handler,
		},
		{
			MethodName: "InitializeBackend",
			Handler:    _Backend_InitializeBackend_Handler,
		},
		{
			MethodName: "ListBackend",
			Handler:    _Backend_ListBackend_Handler,
		},
		{
			MethodName: "TerminateBackend",
			Handler:    _Backend_TerminateBackend_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backend.proto",
}
