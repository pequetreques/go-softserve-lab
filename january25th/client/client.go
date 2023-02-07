package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "softserve.com/january25th/generated/proto"
)

func main() {
	connection, err := grpc.Dial("localhost:7080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln(err)
	}
	defer connection.Close()

	client := pb.NewUsersApiClient(connection)
	res, err := client.GetUser(context.Background(), &pb.IdRequest{Id: "7878"})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res)
}
