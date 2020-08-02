package authentication

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints interface {
	LoginEndpoint() endpoint.Endpoint
}

type endpoints struct {
	service Service
}

func InitEndpoints(service Service) Endpoints {
	return &endpoints{service}
}

func (e *endpoints) LoginEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*LoginRequest)
		token, isRegistered, err := e.service.Login(ctx, req.State, req.Code)
		return LoginResponse{
			Token:        token,
			IsRegistered: isRegistered,
		}, err
	}
}
