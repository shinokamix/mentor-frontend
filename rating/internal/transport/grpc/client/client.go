package grpccleint

import (
	"context"
	"fmt"

	pb "rating/pkg/api/proto"

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

func (m *MentorClient) MethodMentorRating(ctx context.Context, action, mentorEmail string, rating float32) error {
	req := &pb.RatingRequest{
		Action:      action,
		MentorEmail: mentorEmail,
		Rating:      rating,
	}

	resp, err := m.client.MethodMentorRating(ctx, req)
	if err != nil {
		return fmt.Errorf("MethodMentorRating RPC call failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("server responded with success=false, message=%s", resp.Message)
	}

	return nil
}
