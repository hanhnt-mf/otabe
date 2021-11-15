// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// OTabeManagerClient is the client API for OTabeManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OTabeManagerClient interface {
	GetRestaurantDetails(ctx context.Context, in *GetRestaurantRequest, opts ...grpc.CallOption) (*GetRestaurantResponse, error)
	ListRestaurants(ctx context.Context, in *ListRestaurantsRequest, opts ...grpc.CallOption) (*ListRestaurantsResponse, error)
}

type oTabeManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewOTabeManagerClient(cc grpc.ClientConnInterface) OTabeManagerClient {
	return &oTabeManagerClient{cc}
}

func (c *oTabeManagerClient) GetRestaurantDetails(ctx context.Context, in *GetRestaurantRequest, opts ...grpc.CallOption) (*GetRestaurantResponse, error) {
	out := new(GetRestaurantResponse)
	err := c.cc.Invoke(ctx, "/v1.OTabeManager/GetRestaurantDetails", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oTabeManagerClient) ListRestaurants(ctx context.Context, in *ListRestaurantsRequest, opts ...grpc.CallOption) (*ListRestaurantsResponse, error) {
	out := new(ListRestaurantsResponse)
	err := c.cc.Invoke(ctx, "/v1.OTabeManager/ListRestaurants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OTabeManagerServer is the server API for OTabeManager service.
// All implementations must embed UnimplementedOTabeManagerServer
// for forward compatibility
type OTabeManagerServer interface {
	GetRestaurantDetails(context.Context, *GetRestaurantRequest) (*GetRestaurantResponse, error)
	ListRestaurants(context.Context, *ListRestaurantsRequest) (*ListRestaurantsResponse, error)
	mustEmbedUnimplementedOTabeManagerServer()
}

// UnimplementedOTabeManagerServer must be embedded to have forward compatible implementations.
type UnimplementedOTabeManagerServer struct {
}

func (UnimplementedOTabeManagerServer) GetRestaurantDetails(context.Context, *GetRestaurantRequest) (*GetRestaurantResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRestaurantDetails not implemented")
}
func (UnimplementedOTabeManagerServer) ListRestaurants(context.Context, *ListRestaurantsRequest) (*ListRestaurantsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRestaurants not implemented")
}
func (UnimplementedOTabeManagerServer) mustEmbedUnimplementedOTabeManagerServer() {}

// UnsafeOTabeManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OTabeManagerServer will
// result in compilation errors.
type UnsafeOTabeManagerServer interface {
	mustEmbedUnimplementedOTabeManagerServer()
}

func RegisterOTabeManagerServer(s grpc.ServiceRegistrar, srv OTabeManagerServer) {
	s.RegisterService(&OTabeManager_ServiceDesc, srv)
}

func _OTabeManager_GetRestaurantDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRestaurantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OTabeManagerServer).GetRestaurantDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.OTabeManager/GetRestaurantDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OTabeManagerServer).GetRestaurantDetails(ctx, req.(*GetRestaurantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OTabeManager_ListRestaurants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRestaurantsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OTabeManagerServer).ListRestaurants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.OTabeManager/ListRestaurants",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OTabeManagerServer).ListRestaurants(ctx, req.(*ListRestaurantsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OTabeManager_ServiceDesc is the grpc.ServiceDesc for OTabeManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OTabeManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.OTabeManager",
	HandlerType: (*OTabeManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRestaurantDetails",
			Handler:    _OTabeManager_GetRestaurantDetails_Handler,
		},
		{
			MethodName: "ListRestaurants",
			Handler:    _OTabeManager_ListRestaurants_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/otabe.proto",
}