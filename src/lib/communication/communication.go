package communication

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	Collect "../collect"
	pb "../pb/api"
)

// Communication wrapper for gRPC Transport

type DataAPIConnection struct {
	Logger     *log.Logger
	Address    string
	Connection *grpc.ClientConn
	ApiClient  pb.DataApiClient
}

type ApiHandler struct {
	ApiConn *DataAPIConnection
}

func (handler *ApiHandler) SendByteMsg(ctx context.Context, inp *pb.ByteStringRequest) (*pb.ByteStringReply, error) {
	val := inp.ByteString
	handler.ApiConn.Logger.Println(val)

	resp := pb.ByteStringReply{
		Resp: pb.ByteStringReply_SUCCESS,
	}

	return &resp, nil
}

func (handler *ApiHandler) SendJSONMessage(ctx context.Context, inp *pb.ByteStringRequest) (*pb.JSONMessageReply, error) {
	val := inp.ByteString
	handler.ApiConn.Logger.Println(val)

	resp := pb.JSONMessageReply{
		Resp: pb.JSONMessageReply_SUCCESS_JSON,
	}

	return &resp, nil
}

func (handler *ApiHandler) SendSystemInfo(ctx context.Context, inp *pb.SystemInfoRequest) (*pb.SystemInfoReply, error) {
	handler.ApiConn.Logger.Println(inp)

	resp := pb.SystemInfoReply{
		Resp: pb.SystemInfoReply_SUCCESS,
	}

	return &resp, nil
}

func NewConnection(address string, logger *log.Logger) *DataAPIConnection {
	apiConnection := new(DataAPIConnection)
	apiConnection.Address = address
	apiConnection.Logger = logger

	return apiConnection
}

func (apiConn *DataAPIConnection) Dial() {
	// Dial address.
	conn, err := grpc.Dial(apiConn.Address, grpc.WithInsecure())
	if err != nil {
		apiConn.Logger.Fatalf("Could not make connection to server.\nError: %v", err)
	}
	apiConn.Connection = conn

	// Bind as service client.
	apiClient := pb.NewDataApiClient(apiConn.Connection)
	apiConn.ApiClient = apiClient
}

func (apiConn *DataAPIConnection) Listen() {
	lis, err := net.Listen("tcp", apiConn.Address)
	if err != nil {
		apiConn.Logger.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterDataApiServer(server, &ApiHandler{apiConn})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (apiConn DataAPIConnection) SendByteMsg(msg []byte) *pb.ByteStringReply {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	byteMsgRequest := pb.ByteStringRequest{ByteString: msg}
	resp, err := apiConn.ApiClient.SendByteMsg(ctx, &byteMsgRequest)
	if err != nil {
		apiConn.Logger.Fatalf("could not send byte msg: %v", err)
	}

	return resp
}

func (apiConn *DataAPIConnection) SendSystemInfo(sysinfo Collect.SystemInfoStat) *pb.SystemInfoReply {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	systemInfoRequest := pb.SystemInfoRequest{
		TotalMemory:     sysinfo.TotalMemory,
		AvailableMemory: sysinfo.AvailableMemory,
		UsedMemory:      sysinfo.UsedMemory,
		TotalDisk:       sysinfo.TotalDisk,
		FreeDisk:        sysinfo.FreeDisk,
		UsedDisk:        sysinfo.UsedDisk,
		DiskPath:        sysinfo.DiskPath,
		Hostname:        sysinfo.Hostname,
		OS:              sysinfo.OS,
		Timestamp:       sysinfo.Timestamp,
	}

	resp, err := apiConn.ApiClient.SendSystemInfo(ctx, &systemInfoRequest)
	if err != nil {
		apiConn.Logger.Fatalf("could not send system info: %v", err)
	}

	apiConn.Logger.Println(sysinfo)

	return resp
}
