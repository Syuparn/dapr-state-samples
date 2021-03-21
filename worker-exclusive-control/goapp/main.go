package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/xerrors"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

const (
	store = "redis-state"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	s := daprd.NewService(":3000")

	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	if err := s.AddBindingInvocationHandler("/append-message", appendMessageHandler(client)); err != nil {
		log.Fatalf("error adding binding handler: %v", err)
		os.Exit(1)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
}

type MessageAppendedEvent struct {
	EventID   string `json:"event_id"`
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
}

type Message struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func appendMessageHandler(client dapr.Client) func(context.Context, *common.BindingEvent) ([]byte, error) {
	return func(ctx context.Context, in *common.BindingEvent) ([]byte, error) {
		log.Printf("binding - Data:%s", string(in.Data))

		var e MessageAppendedEvent
		if err := json.Unmarshal(in.Data, &e); err != nil {
			return nil, xerrors.Errorf("failed to unmarshal event: %w", err)
		}

		log.Printf("event: %v", e)
		for i, c := range e.Message {
			// get current message
			item, err := client.GetState(ctx, store, e.MessageID)
			if err != nil {
				return nil, xerrors.Errorf("failed to get state (id: %s, message[%d]): %w",
					e.MessageID, i, err)
			}

			var m Message
			if err := json.Unmarshal(item.Value, &m); err != nil {
				return nil, xerrors.Errorf("failed to unmarshal message: %w", err)
			}
			log.Printf("message: %v", m)

			time.Sleep(time.Duration(rand.Int63n(2500)) * time.Millisecond)

			// append one character
			m.Message += string(c)

			log.Printf("updated message: %v", m)
			v, err := json.Marshal(m)
			if err != nil {
				return nil, xerrors.Errorf("failed to marshal message: %w", err)
			}

			// save updated message
			setItem := &dapr.SetStateItem{
				Key:   m.ID,
				Value: v,
				// use etag for optimistic lock!!
				Etag: &dapr.ETag{Value: item.Etag},
				Options: &dapr.StateOptions{
					Concurrency: dapr.StateConcurrencyFirstWrite,
					Consistency: dapr.StateConsistencyStrong,
				},
			}

			if err := client.SaveBulkState(ctx, store, setItem); err != nil {
				log.Printf("failed to save message: %v", m)
				return nil, xerrors.Errorf("failed to save state (id: %s, message[%d]): %w",
					e.MessageID, i, err)
			}

			time.Sleep(time.Duration(rand.Int63n(2500)) * time.Millisecond)
		}

		return nil, nil
	}
}
