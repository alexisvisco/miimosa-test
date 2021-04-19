// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sessions

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

// SessionsClient is the client API for Sessions service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SessionsClient interface {
	// Create a new token
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*TokenReply, error)
	// Validate the token
	Validate(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*TokenReply, error)
}

type sessionsClient struct {
	cc grpc.ClientConnInterface
}

func NewSessionsClient(cc grpc.ClientConnInterface) SessionsClient {
	return &sessionsClient{cc}
}

func (c *sessionsClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*TokenReply, error) {
	out := new(TokenReply)
	err := c.cc.Invoke(ctx, "/sessions.Sessions/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionsClient) Validate(ctx context.Context, in *ValidateTokenRequest, opts ...grpc.CallOption) (*TokenReply, error) {
	out := new(TokenReply)
	err := c.cc.Invoke(ctx, "/sessions.Sessions/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SessionsServer is the server API for Sessions service.
// All implementations must embed UnimplementedSessionsServer
// for forward compatibility
type SessionsServer interface {
	// Create a new token
	Create(context.Context, *CreateRequest) (*TokenReply, error)
	// Validate the token
	Validate(context.Context, *ValidateTokenRequest) (*TokenReply, error)
	mustEmbedUnimplementedSessionsServer()
}

// UnimplementedSessionsServer must be embedded to have forward compatible implementations.
type UnimplementedSessionsServer struct {
}

func (UnimplementedSessionsServer) Create(context.Context, *CreateRequest) (*TokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedSessionsServer) Validate(context.Context, *ValidateTokenRequest) (*TokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedSessionsServer) mustEmbedUnimplementedSessionsServer() {}

// UnsafeSessionsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SessionsServer will
// result in compilation errors.
type UnsafeSessionsServer interface {
	mustEmbedUnimplementedSessionsServer()
}

func RegisterSessionsServer(s grpc.ServiceRegistrar, srv SessionsServer) {
	s.RegisterService(&Sessions_ServiceDesc, srv)
}

func _Sessions_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionsServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sessions.Sessions/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionsServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sessions_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionsServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sessions.Sessions/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionsServer).Validate(ctx, req.(*ValidateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sessions_ServiceDesc is the grpc.ServiceDesc for Sessions service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sessions_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sessions.Sessions",
	HandlerType: (*SessionsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Sessions_Create_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _Sessions_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/sessions.proto",
}
