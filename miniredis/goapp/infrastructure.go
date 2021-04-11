package main

import (
	"context"
	"encoding/json"
	"fmt"

	dapr "github.com/dapr/go-sdk/client"
	"golang.org/x/xerrors"
)

const (
	store = "redis-state"
)

type messageRepository struct {
	c dapr.Client
}

func NewMessageRepository(c dapr.Client) MessageRepository {
	return &messageRepository{
		c: c,
	}
}

func (r *messageRepository) Save(ctx context.Context, message *Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return xerrors.Errorf("failed to marshal message: %w", err)
	}

	if err := r.c.SaveState(ctx, store, message.ID, data); err != nil {
		return xerrors.Errorf("failed to save message: %w", err)
	}
	fmt.Printf("message saved: %s", string(data))

	return nil
}
