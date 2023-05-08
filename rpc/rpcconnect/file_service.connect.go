// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: file_service.proto

package rpcconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	rpc "github.com/gokch/ipfs_mount/rpc"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// FileServiceName is the fully-qualified name of the FileService service.
	FileServiceName = "proto.FileService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// FileServiceUploadProcedure is the fully-qualified name of the FileService's Upload RPC.
	FileServiceUploadProcedure = "/proto.FileService/Upload"
	// FileServiceDownloadProcedure is the fully-qualified name of the FileService's Download RPC.
	FileServiceDownloadProcedure = "/proto.FileService/Download"
)

// FileServiceClient is a client for the proto.FileService service.
type FileServiceClient interface {
	Upload(context.Context, *connect_go.Request[rpc.UploadRequest]) (*connect_go.Response[rpc.UploadResponse], error)
	Download(context.Context, *connect_go.Request[rpc.DownloadRequest]) (*connect_go.Response[rpc.DownloadResponse], error)
}

// NewFileServiceClient constructs a client for the proto.FileService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewFileServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) FileServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &fileServiceClient{
		upload: connect_go.NewClient[rpc.UploadRequest, rpc.UploadResponse](
			httpClient,
			baseURL+FileServiceUploadProcedure,
			opts...,
		),
		download: connect_go.NewClient[rpc.DownloadRequest, rpc.DownloadResponse](
			httpClient,
			baseURL+FileServiceDownloadProcedure,
			opts...,
		),
	}
}

// fileServiceClient implements FileServiceClient.
type fileServiceClient struct {
	upload   *connect_go.Client[rpc.UploadRequest, rpc.UploadResponse]
	download *connect_go.Client[rpc.DownloadRequest, rpc.DownloadResponse]
}

// Upload calls proto.FileService.Upload.
func (c *fileServiceClient) Upload(ctx context.Context, req *connect_go.Request[rpc.UploadRequest]) (*connect_go.Response[rpc.UploadResponse], error) {
	return c.upload.CallUnary(ctx, req)
}

// Download calls proto.FileService.Download.
func (c *fileServiceClient) Download(ctx context.Context, req *connect_go.Request[rpc.DownloadRequest]) (*connect_go.Response[rpc.DownloadResponse], error) {
	return c.download.CallUnary(ctx, req)
}

// FileServiceHandler is an implementation of the proto.FileService service.
type FileServiceHandler interface {
	Upload(context.Context, *connect_go.Request[rpc.UploadRequest]) (*connect_go.Response[rpc.UploadResponse], error)
	Download(context.Context, *connect_go.Request[rpc.DownloadRequest]) (*connect_go.Response[rpc.DownloadResponse], error)
}

// NewFileServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewFileServiceHandler(svc FileServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(FileServiceUploadProcedure, connect_go.NewUnaryHandler(
		FileServiceUploadProcedure,
		svc.Upload,
		opts...,
	))
	mux.Handle(FileServiceDownloadProcedure, connect_go.NewUnaryHandler(
		FileServiceDownloadProcedure,
		svc.Download,
		opts...,
	))
	return "/proto.FileService/", mux
}

// UnimplementedFileServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedFileServiceHandler struct{}

func (UnimplementedFileServiceHandler) Upload(context.Context, *connect_go.Request[rpc.UploadRequest]) (*connect_go.Response[rpc.UploadResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("proto.FileService.Upload is not implemented"))
}

func (UnimplementedFileServiceHandler) Download(context.Context, *connect_go.Request[rpc.DownloadRequest]) (*connect_go.Response[rpc.DownloadResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("proto.FileService.Download is not implemented"))
}
