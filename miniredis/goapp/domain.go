package main

import (
	"context"
	"math/rand"
	"time"

	ulid "github.com/oklog/ulid/v2"
	"golang.org/x/xerrors"
)

type Message struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func NewMessage(msg string) (*Message, error) {
	id, err := newULID()
	if err != nil {
		return nil, xerrors.Errorf("failed to create message: %w", err)
	}

	return &Message{
		ID:      id.String(),
		Message: msg,
	}, nil
}

type MessageRepository interface {
	Save(ctx context.Context, message *Message) error
}

func newULID() (ulid.ULID, error) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.New(ulid.Timestamp(t), entropy)
}
