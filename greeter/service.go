package greeter

import (
	"context"
	"fmt"
	"protobuf-http/proto/pb/greeterpb"

	"github.com/pkg/errors"
)

type serviceBase struct{}

func (s *serviceBase) SayGreeting(ctx context.Context, req *greeterpb.GreeterRequest) (*greeterpb.GreeterResponse, error) {
	if req == nil {
		return nil, errors.New("no parameter")
	}

	if req.Name == "" {
		return nil, errors.New("empty parameter")
	}

	return &greeterpb.GreeterResponse{Greeting: fmt.Sprintf(`hola %s`, req.Name)}, nil
}

func NewService() Service {
	return &serviceBase{}
}
