package main

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/rpcreplay"
	dapr "github.com/dapr/go-sdk/client"
	"google.golang.org/grpc"
)

func TestMessageRepositorySave(t *testing.T) {
	tests := []struct {
		title    string
		cassette string
		msg      *Message
	}{
		{
			"succeed to save",
			"save_success",
			&Message{
				ID:      "0123456789ABCDEFGHJKMNPQRS",
				Message: "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			c, teardown := mockReplayerClient(tt.cassette)
			defer teardown()
			repo := NewMessageRepository(c)

			if err := repo.Save(context.TODO(), tt.msg); err != nil {
				t.Errorf("err must be nil. got=%v", err)
			}
		})
	}
}

func mockRecorderClient(cassette string) (dapr.Client, func()) {
	rec, err := rpcreplay.NewRecorder(
		fmt.Sprintf("testdata/%s.replay", cassette), nil)
	if err != nil {
		panic(err)
	}
	teardown := func() {
		rec.Close()
	}

	conn, err := grpc.Dial("localhost:50001", append(rec.DialOptions(), grpc.WithInsecure())...)
	if err != nil {
		panic(err)
	}

	return dapr.NewClientWithConnection(conn), teardown
}

func mockReplayerClient(cassette string) (dapr.Client, func()) {
	rec, err := rpcreplay.NewReplayer(
		fmt.Sprintf("testdata/%s.replay", cassette))
	if err != nil {
		panic(err)
	}
	teardown := func() {
		rec.Close()
	}

	// HACK: disable rec.DialOptions()[0] (grpc.WithBlock()) otherwise client wait for ping reply
	// from sidecar before replayer starts
	conn, err := grpc.Dial("localhost:50001", append(rec.DialOptions()[1:], grpc.WithInsecure())...)
	if err != nil {
		panic(err)
	}

	return dapr.NewClientWithConnection(conn), teardown
}
