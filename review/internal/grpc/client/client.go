package grpcclient

import (
	"context"
	"fmt"
	pb "review/pkg/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MentorClient struct {
	client pb.MentorServiceClient
	conn   *grpc.ClientConn
}

func NewMentorClient(addr string) (*MentorClient, error) {
	creds := insecure.NewCredentials()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to dial mentor service: %w", err)
	}

	c := pb.NewMentorServiceClient(conn)

	return &MentorClient{
		client: c,
		conn:   conn,
	}, nil
}

func (m *MentorClient) Close() error {
	return m.conn.Close()
}

func (m *MentorClient) CheckMentor(ctx context.Context, mentorEmail string) (bool, error) {
	req := &pb.CheckRequest{
		MentorEmail: mentorEmail,
	}

	resp, err := m.client.CheckMentor(ctx, req)
	if err != nil {
		return false, fmt.Errorf("CheckMentor RPC call failed: %w", err)
	}

	if !resp.Success {
		return false, fmt.Errorf("server responded with failure: %s", resp.Message)
	}

	return resp.Exists, nil
}
