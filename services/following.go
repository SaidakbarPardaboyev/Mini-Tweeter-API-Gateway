package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	userService "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/genproto/users_service"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FollowingServiceI interface {
	FollowingService(ctx context.Context) userService.FollowingServiceClient
}

type FollowingService struct {
	pool *grpcpool.Pool
	conn *grpc.ClientConn
}

func NewFollowingService(cfg *config.Config) (FollowingServiceI, error) {
	factory := func() (*grpc.ClientConn, error) {
		return grpc.NewClient(
			fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800), grpc.MaxCallSendMsgSize(52428800)))
	}

	pool, err := grpcpool.New(factory, 2, 4, time.Minute*30)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800), grpc.MaxCallSendMsgSize(52428800)))

	if err != nil {
		return nil, err
	}

	return &FollowingService{
		pool: pool,
		conn: conn,
	}, nil
}
func (s *FollowingService) FollowingService(ctx context.Context) userService.FollowingServiceClient {
	// conn, err := s.pool.Get(ctx)

	// if err != nil {
	// 	return userService.NewAuthServiceClient(s.conn)
	// }

	return userService.NewFollowingServiceClient(s.conn)
}
