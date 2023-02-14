package Core

import (
	"buf.build/gen/go/kavindudodan/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/kavindudodan/flagd/protocolbuffers/go/sync/v1"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"time"
)

type Server struct {
	Config Config
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", s.Config.Host+":"+s.Config.Port)
	if err != nil {
		log.Printf("Error when listening to address : %s\n", err.Error())
		return
	}

	options, err := s.buildOptions()
	if err != nil {
		log.Printf("Error building dial options : %s\n", err.Error())
		return
	}

	server := grpc.NewServer(options...)
	syncv1grpc.RegisterFlagSyncServiceServer(server, &ServerImpl{})

	fmt.Printf("Server listening : %s\n", hostDefault+":"+portDefault)
	err = server.Serve(listen)
	if err != nil {
		log.Printf("Error when starting the server : %s\n", err.Error())
		return
	}
}

func (s *Server) buildOptions() ([]grpc.ServerOption, error) {
	var options []grpc.ServerOption

	if !s.Config.Secure {
		return options, nil
	}

	keyPair, err := tls.LoadX509KeyPair(s.Config.CertPath, s.Config.KeyPath)
	if err != nil {
		return nil, err
	}

	options = append(options, grpc.Creds(credentials.NewServerTLSFromCert(&keyPair)))

	return options, nil
}

type ServerImpl struct {
}

func (s *ServerImpl) SyncFlags(req *v1.SyncFlagsRequest, stream syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	log.Printf("Requesting flags for : %s", req.ProviderId)

	for _, data := range mockFlagSlice() {
		err := stream.Send(&data)
		if err != nil {
			fmt.Println("Error sending: " + err.Error())
			return err
		}
		time.Sleep(10 * time.Second)
	}

	// long sleep
	for {
		err := stream.Send(&v1.SyncFlagsResponse{
			FlagConfiguration: "",
			State:             v1.SyncState_SYNC_STATE_PING,
		})

		if err != nil {
			fmt.Printf("Error with stream: %s\n", err.Error())
			return err
		}

		time.Sleep(10 * time.Second)
	}
}

func mockFlagSlice() []v1.SyncFlagsResponse {
	return []v1.SyncFlagsResponse{
		{
			FlagConfiguration: fulA,
			State:             v1.SyncState_SYNC_STATE_ALL,
		},
		{
			FlagConfiguration: "",
			State:             v1.SyncState_SYNC_STATE_PING,
		},
		{
			FlagConfiguration: add,
			State:             v1.SyncState_SYNC_STATE_ADD,
		},
		{
			FlagConfiguration: "",
			State:             -1,
		},
		{
			FlagConfiguration: remove,
			State:             v1.SyncState_SYNC_STATE_DELETE,
		},
		{
			FlagConfiguration: fullB,
			State:             v1.SyncState_SYNC_STATE_ALL,
		},
	}
}
