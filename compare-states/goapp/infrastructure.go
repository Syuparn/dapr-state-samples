package main

import (
	"context"
	"encoding/json"

	dapr "github.com/dapr/go-sdk/client"
	"golang.org/x/xerrors"
)

var (
	states = []string{
		"redis-state",
		"mysql-state",
		"mongodb-state",
	}
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

	for _, state := range states {
		if err := r.c.SaveState(ctx, state, message.ID, data); err != nil {
			return xerrors.Errorf("failed to save message to %s: %w", state, err)
		}
	}

	return nil
}
