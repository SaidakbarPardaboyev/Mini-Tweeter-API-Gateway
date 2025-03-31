package services

import (
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
)

type ServiceManager interface {
	UserService() UserServiceI
}

type grpcClients struct {
	userService UserServiceI
}

func NewGrpcClients(conf *config.Config) (ServiceManager, error) {
	userService, err := NewUserService(conf)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		userService: userService,
	}, nil
}

func (g *grpcClients) UserService() UserServiceI {
	return g.userService
}
