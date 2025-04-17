package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	tweetService "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/genproto/tweets_service"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TweetServiceI interface {
	TweetService(ctx context.Context) tweetService.TweetsServiceClient
	TweetMediaService(ctx context.Context) tweetService.TweetMediasServiceClient
}

type TweetService struct {
	pool *grpcpool.Pool
	conn *grpc.ClientConn
}

func NewTweetService(cfg *config.Config) (TweetServiceI, error) {
	factory := func() (*grpc.ClientConn, error) {
		return grpc.NewClient(
			fmt.Sprintf("%s:%d", cfg.TweetServiceHost, cfg.TweetServicePort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800), grpc.MaxCallSendMsgSize(52428800)))
	}

	pool, err := grpcpool.New(factory, 2, 4, time.Minute*30)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.TweetServiceHost, cfg.TweetServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(52428800), grpc.MaxCallSendMsgSize(52428800)))

	if err != nil {
		return nil, err
	}

	return &TweetService{
		pool: pool,
		conn: conn,
	}, nil
}

func (t *TweetService) TweetService(ctx context.Context) tweetService.TweetsServiceClient {
	// conn, err := s.pool.Get(ctx)

	// if err != nil {
	// 	return userService.NewAuthServiceClient(s.conn)
	// }

	return tweetService.NewTweetsServiceClient(t.conn)
}

func (t *TweetService) TweetMediaService(ctx context.Context) tweetService.TweetMediasServiceClient {
	// conn, err := s.pool.Get(ctx)

	// if err != nil {
	// 	return userService.NewAuthServiceClient(s.conn)
	// }

	return tweetService.NewTweetMediasServiceClient(t.conn)
}
