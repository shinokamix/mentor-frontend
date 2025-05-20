package grpcclient

import (
	"context"
	"fmt"
	pb "mentorlink/pkg/api/proto"

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

func (m *MentorClient) NewMentor(ctx context.Context, mentorEmail, contact string) error {
	req := &pb.MentorRequest{
		MentorEmail: mentorEmail,
		Contact:     contact,
	}

	resp, err := m.client.NewMentor(ctx, req)
	if err != nil {
		return fmt.Errorf("NewMentor RPC call failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server resonded with success=false, message=%s", resp.Message)
	}

	return nil
}
