// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package device

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DriverClient is the client API for Driver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DriverClient interface {
	Connect(ctx context.Context, in *ConnectReq, opts ...grpc.CallOption) (*ConnectReply, error)
	ReConnect(ctx context.Context, in *ReConnectReq, opts ...grpc.CallOption) (*ReConnectReply, error)
	DisConnect(ctx context.Context, in *DisConnectReq, opts ...grpc.CallOption) (*DisConnectReply, error)
	ChangeCurrentPage(ctx context.Context, in *ChangeCurrentPageReq, opts ...grpc.CallOption) (*ChangeCurrentPageReply, error)
	Input(ctx context.Context, opts ...grpc.CallOption) (Driver_InputClient, error)
	Output(ctx context.Context, in *OutputReq, opts ...grpc.CallOption) (Driver_OutputClient, error)
}

type driverClient struct {
	cc grpc.ClientConnInterface
}

func NewDriverClient(cc grpc.ClientConnInterface) DriverClient {
	return &driverClient{cc}
}

func (c *driverClient) Connect(ctx context.Context, in *ConnectReq, opts ...grpc.CallOption) (*ConnectReply, error) {
	out := new(ConnectReply)
	err := c.cc.Invoke(ctx, "/device.driver.Driver/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) ReConnect(ctx context.Context, in *ReConnectReq, opts ...grpc.CallOption) (*ReConnectReply, error) {
	out := new(ReConnectReply)
	err := c.cc.Invoke(ctx, "/device.driver.Driver/ReConnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) DisConnect(ctx context.Context, in *DisConnectReq, opts ...grpc.CallOption) (*DisConnectReply, error) {
	out := new(DisConnectReply)
	err := c.cc.Invoke(ctx, "/device.driver.Driver/DisConnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) ChangeCurrentPage(ctx context.Context, in *ChangeCurrentPageReq, opts ...grpc.CallOption) (*ChangeCurrentPageReply, error) {
	out := new(ChangeCurrentPageReply)
	err := c.cc.Invoke(ctx, "/device.driver.Driver/ChangeCurrentPage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverClient) Input(ctx context.Context, opts ...grpc.CallOption) (Driver_InputClient, error) {
	stream, err := c.cc.NewStream(ctx, &Driver_ServiceDesc.Streams[0], "/device.driver.Driver/Input", opts...)
	if err != nil {
		return nil, err
	}
	x := &driverInputClient{stream}
	return x, nil
}

type Driver_InputClient interface {
	Send(*InputReq) error
	CloseAndRecv() (*InputReply, error)
	grpc.ClientStream
}

type driverInputClient struct {
	grpc.ClientStream
}

func (x *driverInputClient) Send(m *InputReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *driverInputClient) CloseAndRecv() (*InputReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(InputReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *driverClient) Output(ctx context.Context, in *OutputReq, opts ...grpc.CallOption) (Driver_OutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &Driver_ServiceDesc.Streams[1], "/device.driver.Driver/Output", opts...)
	if err != nil {
		return nil, err
	}
	x := &driverOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Driver_OutputClient interface {
	Recv() (*OutputReply, error)
	grpc.ClientStream
}

type driverOutputClient struct {
	grpc.ClientStream
}

func (x *driverOutputClient) Recv() (*OutputReply, error) {
	m := new(OutputReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DriverServer is the server API for Driver service.
// All implementations must embed UnimplementedDriverServer
// for forward compatibility
type DriverServer interface {
	Connect(context.Context, *ConnectReq) (*ConnectReply, error)
	ReConnect(context.Context, *ReConnectReq) (*ReConnectReply, error)
	DisConnect(context.Context, *DisConnectReq) (*DisConnectReply, error)
	ChangeCurrentPage(context.Context, *ChangeCurrentPageReq) (*ChangeCurrentPageReply, error)
	Input(Driver_InputServer) error
	Output(*OutputReq, Driver_OutputServer) error
	mustEmbedUnimplementedDriverServer()
}

// UnimplementedDriverServer must be embedded to have forward compatible implementations.
type UnimplementedDriverServer struct {
}

func (UnimplementedDriverServer) Connect(context.Context, *ConnectReq) (*ConnectReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedDriverServer) ReConnect(context.Context, *ReConnectReq) (*ReConnectReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReConnect not implemented")
}
func (UnimplementedDriverServer) DisConnect(context.Context, *DisConnectReq) (*DisConnectReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisConnect not implemented")
}
func (UnimplementedDriverServer) ChangeCurrentPage(context.Context, *ChangeCurrentPageReq) (*ChangeCurrentPageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeCurrentPage not implemented")
}
func (UnimplementedDriverServer) Input(Driver_InputServer) error {
	return status.Errorf(codes.Unimplemented, "method Input not implemented")
}
func (UnimplementedDriverServer) Output(*OutputReq, Driver_OutputServer) error {
	return status.Errorf(codes.Unimplemented, "method Output not implemented")
}
func (UnimplementedDriverServer) mustEmbedUnimplementedDriverServer() {}

// UnsafeDriverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DriverServer will
// result in compilation errors.
type UnsafeDriverServer interface {
	mustEmbedUnimplementedDriverServer()
}

func RegisterDriverServer(s grpc.ServiceRegistrar, srv DriverServer) {
	s.RegisterService(&Driver_ServiceDesc, srv)
}

func _Driver_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/device.driver.Driver/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).Connect(ctx, req.(*ConnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_ReConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReConnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).ReConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/device.driver.Driver/ReConnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).ReConnect(ctx, req.(*ReConnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_DisConnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisConnectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).DisConnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/device.driver.Driver/DisConnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).DisConnect(ctx, req.(*DisConnectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_ChangeCurrentPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeCurrentPageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverServer).ChangeCurrentPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/device.driver.Driver/ChangeCurrentPage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverServer).ChangeCurrentPage(ctx, req.(*ChangeCurrentPageReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Driver_Input_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DriverServer).Input(&driverInputServer{stream})
}

type Driver_InputServer interface {
	SendAndClose(*InputReply) error
	Recv() (*InputReq, error)
	grpc.ServerStream
}

type driverInputServer struct {
	grpc.ServerStream
}

func (x *driverInputServer) SendAndClose(m *InputReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *driverInputServer) Recv() (*InputReq, error) {
	m := new(InputReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Driver_Output_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OutputReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DriverServer).Output(m, &driverOutputServer{stream})
}

type Driver_OutputServer interface {
	Send(*OutputReply) error
	grpc.ServerStream
}

type driverOutputServer struct {
	grpc.ServerStream
}

func (x *driverOutputServer) Send(m *OutputReply) error {
	return x.ServerStream.SendMsg(m)
}

// Driver_ServiceDesc is the grpc.ServiceDesc for Driver service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Driver_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "device.driver.Driver",
	HandlerType: (*DriverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connect",
			Handler:    _Driver_Connect_Handler,
		},
		{
			MethodName: "ReConnect",
			Handler:    _Driver_ReConnect_Handler,
		},
		{
			MethodName: "DisConnect",
			Handler:    _Driver_DisConnect_Handler,
		},
		{
			MethodName: "ChangeCurrentPage",
			Handler:    _Driver_ChangeCurrentPage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Input",
			Handler:       _Driver_Input_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Output",
			Handler:       _Driver_Output_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "driver.proto",
}