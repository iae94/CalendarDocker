package main

import (
	cfg "calendar/pkg/config"
	logging "calendar/pkg/logger"
	pb "calendar/pkg/services/api/gen"
	api "calendar/pkg/services/api/server"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	//Read config
	config, err := cfg.ReadAPIConfig()
	if err != nil {
		log.Fatalf("Reading api config error: %v \n", err)
	}

	logger, err := logging.CreateLogger(&config.API.Logger)
	if err != nil {
		log.Fatalf("Creating api logger error: %v \n", err)
	}

	apiServer, err := api.NewAPIServer(config, logger)
	if err != nil {
		log.Fatalf("Cannot create APIServer instance: %v", err)
	}

	port := fmt.Sprintf(":%v", apiServer.Config.API.Port)

	// TCP Listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC server
	grpcServer := grpc.NewServer()

	// Register our service in grpc server
	pb.RegisterAPIServer(grpcServer, apiServer)

	// Initial log
	apiServer.Logger.Sugar().Infof("Start API service at %v port\n", port)


	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}