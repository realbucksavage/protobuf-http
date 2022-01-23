package greeter

import (
	"context"
	"protobuf-http/proto/pb/greeterpb"
)

type Service interface {
	SayGreeting(ctx context.Context, req *greeterpb.GreeterRequest) (*greeterpb.GreeterResponse, error)
}
