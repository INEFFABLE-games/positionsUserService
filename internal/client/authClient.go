package client

import (
	"context"
	protocol "github.com/INEFFABLE-games/authService/protocol"
)

type AuthClient struct {
	client          protocol.AuthServiceClient
}

func (a *AuthClient)Refresh(ctx context.Context, in *protocol.RefreshRequest) (*protocol.RefreshReply, error){
	reply, err := a.client.Refresh(ctx, &protocol.RefreshRequest{
		Uid: in.Uid,
	})

	return reply, err
}

func NewAuthClient(client protocol.AuthServiceClient) *AuthClient {
	return &AuthClient{client: client}
}