// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package frontend

import (
	context "context"
	message "github.com/afrozalm/minimess/message"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FrontEndClient is the client API for FrontEnd service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrontEndClient interface {
	BroadcastOut(ctx context.Context, in *message.Chat, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type frontEndClient struct {
	cc grpc.ClientConnInterface
}

func NewFrontEndClient(cc grpc.ClientConnInterface) FrontEndClient {
	return &frontEndClient{cc}
}

func (c *frontEndClient) BroadcastOut(ctx context.Context, in *message.Chat, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/frontend.FrontEnd/BroadcastOut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FrontEndServer is the server API for FrontEnd service.
// All implementations must embed UnimplementedFrontEndServer
// for forward compatibility
type FrontEndServer interface {
	BroadcastOut(context.Context, *message.Chat) (*emptypb.Empty, error)
	mustEmbedUnimplementedFrontEndServer()
}

// UnimplementedFrontEndServer must be embedded to have forward compatible implementations.
type UnimplementedFrontEndServer struct {
}

func (UnimplementedFrontEndServer) BroadcastOut(context.Context, *message.Chat) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastOut not implemented")
}
func (UnimplementedFrontEndServer) mustEmbedUnimplementedFrontEndServer() {}

// UnsafeFrontEndServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrontEndServer will
// result in compilation errors.
type UnsafeFrontEndServer interface {
	mustEmbedUnimplementedFrontEndServer()
}

func RegisterFrontEndServer(s grpc.ServiceRegistrar, srv FrontEndServer) {
	s.RegisterService(&FrontEnd_ServiceDesc, srv)
}

func _FrontEnd_BroadcastOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(message.Chat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontEndServer).BroadcastOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/frontend.FrontEnd/BroadcastOut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontEndServer).BroadcastOut(ctx, req.(*message.Chat))
	}
	return interceptor(ctx, in, info, handler)
}

// FrontEnd_ServiceDesc is the grpc.ServiceDesc for FrontEnd service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FrontEnd_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "frontend.FrontEnd",
	HandlerType: (*FrontEndServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BroadcastOut",
			Handler:    _FrontEnd_BroadcastOut_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "frontend.proto",
}
