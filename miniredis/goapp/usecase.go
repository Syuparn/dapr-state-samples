package main

import (
	"context"

	"golang.org/x/xerrors"
)

type CreateMessageInputPort interface {
	Do(ctx context.Context, in *CreateMessageInputData) (*CreateMessageOutputData, error)
}

type CreateMessageInputData struct {
	Message string
}

type CreateMessageOutputData struct {
	Message *Message
}

func NewCreateMessagePort(messageRepository MessageRepository) CreateMessageInputPort {
	return &createMessageInteractor{
		messageRepository: messageRepository,
	}
}

type createMessageInteractor struct {
	messageRepository MessageRepository
}

func (r *createMessageInteractor) Do(ctx context.Context, in *CreateMessageInputData) (*CreateMessageOutputData, error) {
	msg, err := NewMessage(in.Message)
	if err != nil {
		return nil, xerrors.Errorf("failed to create message: %w", err)
	}

	if err := r.messageRepository.Save(ctx, msg); err != nil {
		return nil, xerrors.Errorf("failed to save message: %w", err)
	}

	return &CreateMessageOutputData{
		Message: msg,
	}, nil
}
