// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package fanout

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FanoutClient is the client API for Fanout service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FanoutClient interface {
	Subscribe(ctx context.Context, in *HostTopic, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Unsubscribe(ctx context.Context, in *HostTopic, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnsubscribeAll(ctx context.Context, in *Host, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type fanoutClient struct {
	cc grpc.ClientConnInterface
}

func NewFanoutClient(cc grpc.ClientConnInterface) FanoutClient {
	return &fanoutClient{cc}
}

func (c *fanoutClient) Subscribe(ctx context.Context, in *HostTopic, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/fanout.Fanout/Subscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fanoutClient) Unsubscribe(ctx context.Context, in *HostTopic, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/fanout.Fanout/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fanoutClient) UnsubscribeAll(ctx context.Context, in *Host, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/fanout.Fanout/UnsubscribeAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FanoutServer is the server API for Fanout service.
// All implementations must embed UnimplementedFanoutServer
// for forward compatibility
type FanoutServer interface {
	Subscribe(context.Context, *HostTopic) (*emptypb.Empty, error)
	Unsubscribe(context.Context, *HostTopic) (*emptypb.Empty, error)
	UnsubscribeAll(context.Context, *Host) (*emptypb.Empty, error)
	mustEmbedUnimplementedFanoutServer()
}

// UnimplementedFanoutServer must be embedded to have forward compatible implementations.
type UnimplementedFanoutServer struct {
}

func (UnimplementedFanoutServer) Subscribe(context.Context, *HostTopic) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedFanoutServer) Unsubscribe(context.Context, *HostTopic) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribe not implemented")
}
func (UnimplementedFanoutServer) UnsubscribeAll(context.Context, *Host) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeAll not implemented")
}
func (UnimplementedFanoutServer) mustEmbedUnimplementedFanoutServer() {}

// UnsafeFanoutServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FanoutServer will
// result in compilation errors.
type UnsafeFanoutServer interface {
	mustEmbedUnimplementedFanoutServer()
}

func RegisterFanoutServer(s grpc.ServiceRegistrar, srv FanoutServer) {
	s.RegisterService(&Fanout_ServiceDesc, srv)
}

func _Fanout_Subscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostTopic)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FanoutServer).Subscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fanout.Fanout/Subscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FanoutServer).Subscribe(ctx, req.(*HostTopic))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fanout_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostTopic)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FanoutServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fanout.Fanout/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FanoutServer).Unsubscribe(ctx, req.(*HostTopic))
	}
	return interceptor(ctx, in, info, handler)
}

func _Fanout_UnsubscribeAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Host)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FanoutServer).UnsubscribeAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fanout.Fanout/UnsubscribeAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FanoutServer).UnsubscribeAll(ctx, req.(*Host))
	}
	return interceptor(ctx, in, info, handler)
}

// Fanout_ServiceDesc is the grpc.ServiceDesc for Fanout service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fanout_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fanout.Fanout",
	HandlerType: (*FanoutServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Subscribe",
			Handler:    _Fanout_Subscribe_Handler,
		},
		{
			MethodName: "Unsubscribe",
			Handler:    _Fanout_Unsubscribe_Handler,
		},
		{
			MethodName: "UnsubscribeAll",
			Handler:    _Fanout_UnsubscribeAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fanout.proto",
}
