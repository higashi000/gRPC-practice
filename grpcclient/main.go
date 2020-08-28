package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/higashi000/practice_pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:19003", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()

	ctx := context.Background()
	client := practice_pb.NewStreamTestClient(conn)

	stream, err := client.Test(ctx)
	if err != nil {
		log.Println(err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				fmt.Println(err)
			}

			log.Printf("%s: %s\n", in.Name, in.Text)
		}
	}()

	if err := stream.Send(&practice_pb.TestRequest{Name: "higashi", Text: "Hi!"}); err != nil {
		log.Fatalf("Failed to send a message: %v", err)
	}

	stream.CloseSend()
	<-waitc
}
