// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.4
// source: gserver/invoke/api/register.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RegisterRouter_RegisterRouter_FullMethodName = "/invoke.RegisterRouter/RegisterRouter"
)

// RegisterRouterClient is the client API for RegisterRouter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterRouterClient interface {
	RegisterRouter(ctx context.Context, in *RegisterRouterReq, opts ...grpc.CallOption) (*RegisterRouterResp, error)
}

type registerRouterClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterRouterClient(cc grpc.ClientConnInterface) RegisterRouterClient {
	return &registerRouterClient{cc}
}

func (c *registerRouterClient) RegisterRouter(ctx context.Context, in *RegisterRouterReq, opts ...grpc.CallOption) (*RegisterRouterResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterRouterResp)
	err := c.cc.Invoke(ctx, RegisterRouter_RegisterRouter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterRouterServer is the server API for RegisterRouter service.
// All implementations must embed UnimplementedRegisterRouterServer
// for forward compatibility.
type RegisterRouterServer interface {
	RegisterRouter(context.Context, *RegisterRouterReq) (*RegisterRouterResp, error)
	mustEmbedUnimplementedRegisterRouterServer()
}

// UnimplementedRegisterRouterServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRegisterRouterServer struct{}

func (UnimplementedRegisterRouterServer) RegisterRouter(context.Context, *RegisterRouterReq) (*RegisterRouterResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterRouter not implemented")
}
func (UnimplementedRegisterRouterServer) mustEmbedUnimplementedRegisterRouterServer() {}
func (UnimplementedRegisterRouterServer) testEmbeddedByValue()                        {}

// UnsafeRegisterRouterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterRouterServer will
// result in compilation errors.
type UnsafeRegisterRouterServer interface {
	mustEmbedUnimplementedRegisterRouterServer()
}

func RegisterRegisterRouterServer(s grpc.ServiceRegistrar, srv RegisterRouterServer) {
	// If the following call pancis, it indicates UnimplementedRegisterRouterServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RegisterRouter_ServiceDesc, srv)
}

func _RegisterRouter_RegisterRouter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRouterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterRouterServer).RegisterRouter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegisterRouter_RegisterRouter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterRouterServer).RegisterRouter(ctx, req.(*RegisterRouterReq))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterRouter_ServiceDesc is the grpc.ServiceDesc for RegisterRouter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterRouter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "invoke.RegisterRouter",
	HandlerType: (*RegisterRouterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterRouter",
			Handler:    _RegisterRouter_RegisterRouter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gserver/invoke/api/register.proto",
}
