package services

import (
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
)

type ServiceManager interface {
	UserService() UserServiceI
	FollowingService() FollowingServiceI
	TweetService() TweetServiceI
}

type grpcClients struct {
	userService      UserServiceI
	followingService FollowingServiceI
	tweetService     TweetServiceI
}

func NewGrpcClients(conf *config.Config) (ServiceManager, error) {
	userService, err := NewUserService(conf)
	if err != nil {
		return nil, err
	}

	followingService, err := NewFollowingService(conf)
	if err != nil {
		return nil, err
	}

	tweetService, err := NewTweetService(conf)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		userService:      userService,
		followingService: followingService,
		tweetService:     tweetService,
	}, nil
}

func (g *grpcClients) UserService() UserServiceI {
	return g.userService
}

func (g *grpcClients) FollowingService() FollowingServiceI {
	return g.followingService
}

func (g *grpcClients) TweetService() TweetServiceI {
	return g.tweetService
}
