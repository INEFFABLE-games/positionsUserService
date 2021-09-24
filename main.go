package main

import (
	"database/sql"
	"fmt"
	authService "github.com/INEFFABLE-games/authService/protocol"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"positionsUserService/internal/client"
	"positionsUserService/internal/config"
	"positionsUserService/internal/repository"
	"positionsUserService/internal/server"
	"positionsUserService/internal/service"
	"positionsUserService/protocol"
)

func main() {

	cfg := config.NewConfig()

	conn, err := grpc.Dial(fmt.Sprintf(":%s", cfg.AuthPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	authServ := authService.NewAuthServiceClient(conn)

	authClient := client.NewAuthClient(authServ)

	sqlConn, err := sql.Open("postgres", cfg.PsqlURI)
	if err != nil {
		log.Errorf("unable to connect with postgres %v", err)
	}

	grpcServer := grpc.NewServer()
	userRepository := repository.NewUserRepository(sqlConn)
	userService := service.NewUserService(userRepository,authClient)
	userServer := server.NewUserServer(userService)
	protocol.RegisterUserServiceServer(grpcServer, userServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Errorf("unable to start grpc server %v", err.Error())
	}

}
