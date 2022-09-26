package main

import (
	"context"
	"github.com/diegokrule/gRpc/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	client()
}


func client(){
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c:=proto.NewChatServiceClient(conn)
	ctx:=context.Background()
	md:=metadata.MD{}
	md.Set("custom","value")
	md.Set("c","v")
	ctx1:=metadata.NewOutgoingContext(ctx,md)
	response, err := c.SayHello(ctx1, &proto.Message{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)

}