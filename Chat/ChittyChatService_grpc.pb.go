// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package Chat

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

// ChittyChatServiceClient is the client API for ChittyChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChittyChatServiceClient interface {
	JoinChat(ctx context.Context, in *User, opts ...grpc.CallOption) (ChittyChatService_JoinChatClient, error)
	Publish(ctx context.Context, in *FromClient, opts ...grpc.CallOption) (*Empty, error)
	LeaveChat(ctx context.Context, in *User, opts ...grpc.CallOption) (*Empty, error)
}

type chittyChatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChittyChatServiceClient(cc grpc.ClientConnInterface) ChittyChatServiceClient {
	return &chittyChatServiceClient{cc}
}

func (c *chittyChatServiceClient) JoinChat(ctx context.Context, in *User, opts ...grpc.CallOption) (ChittyChatService_JoinChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChittyChatService_ServiceDesc.Streams[0], "/Chat.ChittyChatService/JoinChat", opts...)
	if err != nil {
		return nil, err
	}
	x := &chittyChatServiceJoinChatClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChittyChatService_JoinChatClient interface {
	Recv() (*FromServer, error)
	grpc.ClientStream
}

type chittyChatServiceJoinChatClient struct {
	grpc.ClientStream
}

func (x *chittyChatServiceJoinChatClient) Recv() (*FromServer, error) {
	m := new(FromServer)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chittyChatServiceClient) Publish(ctx context.Context, in *FromClient, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Chat.ChittyChatService/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chittyChatServiceClient) LeaveChat(ctx context.Context, in *User, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Chat.ChittyChatService/LeaveChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChittyChatServiceServer is the server API for ChittyChatService service.
// All implementations must embed UnimplementedChittyChatServiceServer
// for forward compatibility
type ChittyChatServiceServer interface {
	JoinChat(*User, ChittyChatService_JoinChatServer) error
	Publish(context.Context, *FromClient) (*Empty, error)
	LeaveChat(context.Context, *User) (*Empty, error)
	mustEmbedUnimplementedChittyChatServiceServer()
}

// UnimplementedChittyChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChittyChatServiceServer struct {
}

func (UnimplementedChittyChatServiceServer) JoinChat(*User, ChittyChatService_JoinChatServer) error {
	return status.Errorf(codes.Unimplemented, "method JoinChat not implemented")
}
func (UnimplementedChittyChatServiceServer) Publish(context.Context, *FromClient) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedChittyChatServiceServer) LeaveChat(context.Context, *User) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveChat not implemented")
}
func (UnimplementedChittyChatServiceServer) mustEmbedUnimplementedChittyChatServiceServer() {}

// UnsafeChittyChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChittyChatServiceServer will
// result in compilation errors.
type UnsafeChittyChatServiceServer interface {
	mustEmbedUnimplementedChittyChatServiceServer()
}

func RegisterChittyChatServiceServer(s grpc.ServiceRegistrar, srv ChittyChatServiceServer) {
	s.RegisterService(&ChittyChatService_ServiceDesc, srv)
}

func _ChittyChatService_JoinChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(User)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChittyChatServiceServer).JoinChat(m, &chittyChatServiceJoinChatServer{stream})
}

type ChittyChatService_JoinChatServer interface {
	Send(*FromServer) error
	grpc.ServerStream
}

type chittyChatServiceJoinChatServer struct {
	grpc.ServerStream
}

func (x *chittyChatServiceJoinChatServer) Send(m *FromServer) error {
	return x.ServerStream.SendMsg(m)
}

func _ChittyChatService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FromClient)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Chat.ChittyChatService/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServiceServer).Publish(ctx, req.(*FromClient))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChittyChatService_LeaveChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChittyChatServiceServer).LeaveChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Chat.ChittyChatService/LeaveChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChittyChatServiceServer).LeaveChat(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

// ChittyChatService_ServiceDesc is the grpc.ServiceDesc for ChittyChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChittyChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Chat.ChittyChatService",
	HandlerType: (*ChittyChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _ChittyChatService_Publish_Handler,
		},
		{
			MethodName: "LeaveChat",
			Handler:    _ChittyChatService_LeaveChat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "JoinChat",
			Handler:       _ChittyChatService_JoinChat_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "Chat/ChittyChatService.proto",
}
