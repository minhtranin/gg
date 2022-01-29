package main

import (
	"context"
	"fmt"
	__ "gg/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := __.NewGreetServiceClient(conn)
	// fmt.Printf("created client: %f", c)
	req := __.GreetRequest{
		Greeting: &__.Greeting{
			FirstName: "Minh",
			LastName:  "PRO",
		},
	}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling greet RPC: %v", err)
	}
	log.Printf("response from greet: %v", res.Result)
	// doStreamingGreet(c)
	doStreamingServer(c)
}

func doStreamingGreet(c __.GreetServiceClient) {
	fmt.Println("start streaming greet")
	req := __.GreetManyTimesRequest{
		Greeting: &__.Greeting{
			FirstName: "Minhs",
			LastName:  "Pros",
		},
	}
	reqStream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling greet many times: %v", err)
	}
	for {
		msg, err := reqStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("response from greet many times: %v", msg.GetResult())
	}
}

func doStreamingServer(c __.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	req := []*__.LongGreetRequest{
		&__.LongGreetRequest{
			Greeting: &__.Greeting{
				FirstName: "sMinh",
				LastName:  "sPro",
			},
		},
		&__.LongGreetRequest{
			Greeting: &__.Greeting{
				FirstName: "sMinh",
				LastName:  "sVip",
			},
		},
		&__.LongGreetRequest{
			Greeting: &__.Greeting{
				FirstName: "sMinh",
				LastName:  "s1",
			},
		},
	}
	if err != nil {
		log.Fatalf("error while calling stream server: %v", err)
	}
	for _, r := range req {
		fmt.Printf("sending req: %v\n", r)
		stream.Send(r)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from long great: %v", err)
	}
	fmt.Printf("long greet response: %v", res)
}
