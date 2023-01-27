package main

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"
	pb "softserve.com/january25th/generated/proto"

	"google.golang.org/grpc"
)

type usersApiServer struct {
	pb.UnimplementedUsersApiServer
}

func (s *usersApiServer) GetUsers(ctx context.Context, request *pb.Id) (*pb.User, error) {
	// var user pb.User
	// user.Id = uuid.NewString()
	// user.Name = "Edgar"
	// user.Surname = "Fuentes"
	// user.Age = 47

	return &pb.User{Id: uuid.NewString(), Name: "Edgar", Surname: "Fuentes", Age: 47}, nil
}

func main() {
	listener, err := net.Listen("tcp", "localhost:7080")

	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()
	pb.RegisterUsersApiServer(server, &usersApiServer{})

	err = server.Serve(listener)

	if err != nil {
		log.Fatalln(err)
	}
}
