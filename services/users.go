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

type UserServiceI interface {
	UserService(ctx context.Context) userService.UserServiceClient
}

type UserService struct {
	pool *grpcpool.Pool
	conn *grpc.ClientConn
}

func NewUserService(cfg *config.Config) (UserServiceI, error) {
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

	return &UserService{
		pool: pool,
		conn: conn,
	}, nil
}
func (s *UserService) UserService(ctx context.Context) userService.UserServiceClient {
	// conn, err := s.pool.Get(ctx)

	// if err != nil {
	// 	return userService.NewAuthServiceClient(s.conn)
	// }

	return userService.NewUserServiceClient(s.conn)
}
