package main

import (
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "user/proto"
)

// TestGetDownURL 测试生产下载连接
func TestMain(t *testing.T) {
	creds, err := credentials.NewClientTLSFromFile("certs/server.pem", "localhost")
	if err != nil {
		t.Log("Failed to create TLS credentials :")
		return
	}
	conn, err := grpc.Dial(":8880", grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Log("conn Failed :")
		return
	}
	defer conn.Close()

	c := pb.NewUserClient(conn)
	context := context.Background()
	body := &pb.UserRequest{
		Name:     "admin",
		Password: 123,
		Code:     123,
		Codeid:   123,
	}

	r, err := c.Louder(context, body)
	if err != nil {
		t.Log("Louder Failed :")
		return
	}

	t.Log("Success :", r.JWT)
}
