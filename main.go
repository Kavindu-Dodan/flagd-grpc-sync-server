package main

import (
	"buf.build/gen/go/kavindudodan/flagd/grpc/go/sync/v1/syncv1grpc"
	v1 "buf.build/gen/go/kavindudodan/flagd/protocolbuffers/go/sync/v1"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

const host = "localhost"
const port = "9090"

func main() {
	listen, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Printf("Error when listening to address : %s\n", err.Error())
		return
	}

	server := grpc.NewServer()
	syncv1grpc.RegisterFlagSyncServiceServer(server, &ServerImpl{})

	fmt.Printf("Server listening : %s\n", host+":"+port)
	err = server.Serve(listen)
	if err != nil {
		log.Printf("Error when starting the server : %s\n", err.Error())
		return
	}
}

type ServerImpl struct {
}

func (s *ServerImpl) SyncFlags(req *v1.SyncFlagsRequest, stream syncv1grpc.FlagSyncService_SyncFlagsServer) error {
	log.Printf("Requesting flags for : %s", req.ProviderId)

	for _, data := range gemFlagSlice() {
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

func gemFlagSlice() []v1.SyncFlagsResponse {
	return []v1.SyncFlagsResponse{
		{
			FlagConfiguration: readJson("flags/full.json"),
			State:             v1.SyncState_SYNC_STATE_ALL,
		},
		{
			FlagConfiguration: "",
			State:             v1.SyncState_SYNC_STATE_PING,
		},
		{
			FlagConfiguration: readJson("flags/add.json"),
			State:             v1.SyncState_SYNC_STATE_ADD,
		},
		{
			FlagConfiguration: "",
			State:             41,
		},
		{
			FlagConfiguration: readJson("flags/remove.json"),
			State:             v1.SyncState_SYNC_STATE_DELETE,
		},
		{
			FlagConfiguration: readJson("flags/full2.json"),
			State:             v1.SyncState_SYNC_STATE_ALL,
		},
	}
}

func readJson(file string) string {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Error reading bytes: %s", err)
	}

	return string(bytes)
}
