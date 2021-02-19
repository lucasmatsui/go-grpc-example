package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/lucasmatsui/go-grpc-example/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	//AddUser(client)

	//AddUserVerbose(client)

	//AddUsers(client)

	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Marcio",
		Email: "marcio@marcio.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Marcio",
		Email: "marcio@marcio.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		fmt.Println("Status:", stream.GetStatus())

		if stream.GetStatus() == "User has been inserted" {
			fmt.Println(stream.GetUser())
		}
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "k1",
			Name:  "Kenzo",
			Email: "ken@ken.com",
		},
		{
			Id:    "k2",
			Name:  "Kenzo 2",
			Email: "ken2@ken.com",
		},
		{
			Id:    "k3",
			Name:  "Kenzo 3",
			Email: "ken3@ken.com",
		},
		{
			Id:    "k4",
			Name:  "Kenzo 4",
			Email: "ken4@ken.com",
		},
		{
			Id:    "k5",
			Name:  "Kenzo 5",
			Email: "ken5@ken.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error for receive response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "k1",
			Name:  "Kenzo",
			Email: "ken@ken.com",
		},
		{
			Id:    "k2",
			Name:  "Kenzo 2",
			Email: "ken2@ken.com",
		},
		{
			Id:    "k3",
			Name:  "Kenzo 3",
			Email: "ken3@ken.com",
		},
		{
			Id:    "k4",
			Name:  "Kenzo 4",
			Email: "ken4@ken.com",
		},
		{
			Id:    "k5",
			Name:  "Kenzo 5",
			Email: "ken5@ken.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}

			fmt.Printf("Recebendo user %v com status: %v \n", res.User.GetName(), res.GetStatus())
		}

		close(wait)
	}()

	<-wait
}
